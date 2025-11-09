package model

import "time"

type Document struct {
	ID               uint      `json:"id"`
	TemplateName     string    `json:"template_name"`
	Format           string    `json:"format"`
	Content          string    `json:"content"`
	IsPremium        bool      `json:"is_premium"`
	WatermarkApplied bool      `json:"watermark_applied"`
	Price            float64   `json:"price"`
	CreatedAt        time.Time `json:"created_at"`
}

type TemplateField struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Options     []string `json:"options,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Template struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Content     string          `json:"-"`
	Fields      []TemplateField `json:"fields"`
	Tags        []string        `json:"tags"`
}

type GenerateRequest struct {
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
	Format   string            `json:"format"`
	Premium  bool              `json:"premium"`
}

type Account struct {
	Balance  float64   `json:"balance"`
	Currency string    `json:"currency"`
	Updated  time.Time `json:"updated_at"`
}

type Transaction struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Balance   float64   `json:"balance"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
