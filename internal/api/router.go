package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"law_for_it/internal/document"
	"law_for_it/internal/errs"
	"law_for_it/internal/model"
	"law_for_it/internal/payment"
)

type Server struct {
	docs     *document.Service
	payments *payment.Service
	info     ServerInfo
	mux      *http.ServeMux
}

type ServerInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

func New(docs *document.Service, payments *payment.Service, info ServerInfo) *Server {
	srv := &Server{docs: docs, payments: payments, info: info, mux: http.NewServeMux()}
	srv.registerRoutes()
	return srv
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/status", s.handleStatus)
	s.mux.HandleFunc("/templates", s.handleTemplates)
	s.mux.HandleFunc("/generate", s.handleGenerate)
	s.mux.HandleFunc("/downloadf/", s.handleDownloadFree)
	s.mux.HandleFunc("/downloadp/", s.handleDownloadPremium)
	s.mux.HandleFunc("/balance", s.handleBalance)
	s.mux.HandleFunc("/payhistory", s.handlePayHistory)
	s.mux.HandleFunc("/topup", s.handleTopUp)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"name":    s.info.Name,
		"version": s.info.Version,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleTemplates(w http.ResponseWriter, r *http.Request) {
	templates := s.docs.ListTemplates()
	writeJSON(w, http.StatusOK, map[string]interface{}{"templates": templates})
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req model.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	doc, err := s.docs.Generate(r.Context(), req)
	if err != nil {
		switch err {
		case document.ErrTemplateNotFound:
			writeError(w, http.StatusNotFound, err.Error())
			return
		case document.ErrUnsupportedFormat:
			writeError(w, http.StatusBadRequest, err.Error())
			return
		case document.ErrInsufficientBalance:
			writeError(w, http.StatusPaymentRequired, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":            doc.ID,
		"format":        doc.Format,
		"premium":       doc.IsPremium,
		"watermark":     doc.WatermarkApplied,
		"price":         doc.Price,
		"download_free": fmt.Sprintf("/downloadf/%d", doc.ID),
		"download_full": fmt.Sprintf("/downloadp/%d", doc.ID),
		"generated_at":  doc.CreatedAt,
	})
}

func (s *Server) handleDownloadFree(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/downloadf/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid document id")
		return
	}
	doc, err := s.docs.GetDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, document.ErrDocumentNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(s.docs.WatermarkedPreview(doc)))
}

func (s *Server) handleDownloadPremium(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/downloadp/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid document id")
		return
	}
	doc, err := s.docs.GetDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, document.ErrDocumentNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if !doc.IsPremium {
		writeError(w, http.StatusForbidden, "premium version not purchased")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(doc.Content))
}

func (s *Server) handleBalance(w http.ResponseWriter, r *http.Request) {
	account, err := s.payments.Balance(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (s *Server) handlePayHistory(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	history, err := s.payments.History(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"transactions": history})
}

func (s *Server) handleTopUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Amount  float64 `json:"amount"`
		Message string  `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	account, err := s.payments.TopUp(r.Context(), req.Amount, req.Message)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidAmount) {
			writeError(w, http.StatusBadRequest, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func parseID(path, prefix string) (uint, error) {
	if !strings.HasPrefix(path, prefix) {
		return 0, fmt.Errorf("invalid path")
	}
	idStr := strings.TrimPrefix(path, prefix)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
