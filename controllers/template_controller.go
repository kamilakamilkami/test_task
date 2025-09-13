package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"project/internal/domain"
	"project/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	// "time"
	"strconv"
	"encoding/base64"
	"github.com/aymerick/raymond"
)

type TemplateHandler struct {
	service service.TemplateService
}

func NewTemplateHandler(s service.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: s}
}

// POST /Templates
func (h *TemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var template domain.Template
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newtemplate, err := h.service.Create(context.Background(), &template)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newtemplate)
}

// GET /Templates
func (h *TemplateHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    q := r.URL.Query()

    var (
        templateType *string
        active       *bool
    )

    if v := q.Get("type"); v != "" {
        templateType = &v
    }
    if v := q.Get("active"); v != "" {
        if b, err := strconv.ParseBool(v); err == nil {
            active = &b
        }
    }

    templates, err := h.service.GetAll(ctx, templateType, active)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(templates)
}



// GET /Templates/{id}
func (h *TemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	template, err := h.service.GetByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(template)
}

// PUT /Templates/{id}
func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var template domain.Template
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newtemplate, err := h.service.Update(context.Background(), id, &template)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(newtemplate)
}

// POST /Templates/{id}/Preview
func (h *TemplateHandler) Preview(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := uuid.Parse(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    tmpl, err := h.service.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var req domain.PreviewRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    html, err := raymond.Render(tmpl.Body, req.Data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    pdfBytes, err := h.service.FromHtmlToPdf(r.Context(), html)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    encoded := base64.StdEncoding.EncodeToString(pdfBytes)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "pdfBase64": encoded,
    })
}
