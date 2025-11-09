package document

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"law_for_it/internal/config"
	"law_for_it/internal/errs"
	"law_for_it/internal/model"
	"law_for_it/internal/state"
)

var (
	ErrTemplateNotFound    = errors.New("template not found")
	ErrUnsupportedFormat   = errors.New("unsupported format")
	ErrInsufficientBalance = errs.ErrInsufficientBalance
	ErrDocumentNotFound    = errors.New("document not found")
)

type PremiumCharger interface {
	Charge(ctx context.Context, amount float64, message string) error
}

type Service struct {
	manager      *state.Manager
	templates    map[string]model.Template
	allowedForms map[string]struct{}
	watermark    string
	premiumPrice float64
	charger      PremiumCharger
}

func NewService(manager *state.Manager, cfg config.DocumentConfig, charger PremiumCharger) *Service {
	allowed := make(map[string]struct{}, len(cfg.Formats))
	for _, format := range cfg.Formats {
		allowed[strings.ToLower(strings.TrimSpace(format))] = struct{}{}
	}
	tmplMap := make(map[string]model.Template, len(defaultTemplates))
	for _, tmpl := range defaultTemplates {
		tmplMap[tmpl.Name] = tmpl
	}
	return &Service{
		manager:      manager,
		templates:    tmplMap,
		allowedForms: allowed,
		watermark:    cfg.Watermark,
		premiumPrice: cfg.PremiumPrice,
		charger:      charger,
	}
}

func (s *Service) ListTemplates() []model.Template {
	templates := make([]model.Template, 0, len(s.templates))
	for _, tmpl := range s.templates {
		templates = append(templates, tmpl)
	}
	return templates
}

func (s *Service) WatermarkedPreview(doc *model.Document) string {
	if doc == nil {
		return ""
	}
	if doc.WatermarkApplied {
		return doc.Content
	}
	return s.applyWatermark(doc.Content)
}

func (s *Service) Generate(ctx context.Context, req model.GenerateRequest) (*model.Document, error) {
	templateName := strings.TrimSpace(req.Template)
	tmpl, ok := s.templates[templateName]
	if !ok {
		return nil, ErrTemplateNotFound
	}

	format := strings.ToLower(strings.TrimSpace(req.Format))
	if format == "" {
		format = "pdf"
	}
	if _, supported := s.allowedForms[format]; !supported {
		return nil, ErrUnsupportedFormat
	}

	content, err := s.renderTemplate(tmpl.Content, req.Data)
	if err != nil {
		return nil, err
	}

	doc := &model.Document{
		TemplateName: tmpl.Name,
		Format:       format,
		Content:      content,
		IsPremium:    req.Premium,
		Price:        0,
		CreatedAt:    time.Now().UTC(),
	}

	if req.Premium {
		if s.charger != nil {
			if err := s.charger.Charge(ctx, s.premiumPrice, fmt.Sprintf("Premium document: %s", tmpl.DisplayName)); err != nil {
				if errors.Is(err, ErrInsufficientBalance) {
					return nil, ErrInsufficientBalance
				}
				return nil, err
			}
		}
		doc.Price = s.premiumPrice
	} else {
		doc.Content = s.applyWatermark(doc.Content)
		doc.WatermarkApplied = true
	}

	if err := s.manager.Update(func(data *state.Data) error {
		doc.ID = data.NextDocumentID
		data.NextDocumentID++
		data.Documents = append(data.Documents, *doc)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to store document: %w", err)
	}

	return doc, nil
}

func (s *Service) GetDocument(ctx context.Context, id uint) (*model.Document, error) {
	var document *model.Document
	err := s.manager.View(func(data *state.Data) error {
		for i := range data.Documents {
			if data.Documents[i].ID == id {
				copy := data.Documents[i]
				document = &copy
				return nil
			}
		}
		return ErrDocumentNotFound
	})
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, ErrDocumentNotFound
	}
	return document, nil
}

func (s *Service) renderTemplate(content string, data map[string]string) (string, error) {
	tpl, err := template.New("document").Funcs(template.FuncMap{"uppercase": strings.ToUpper}).Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	buf := bytes.Buffer{}
	if data == nil {
		data = map[string]string{}
	}
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

func (s *Service) applyWatermark(content string) string {
	if s.watermark == "" {
		return content
	}
	return fmt.Sprintf("%s\n\n---\n%s", content, s.watermark)
}
