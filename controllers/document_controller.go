package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"project/internal/domain"
	"project/internal/service"

	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"project/middlewares"
	"html/template"
)

type DocumentHandler struct {
	service service.DocumentService
}

func NewDocumentHandler(s service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: s}
}

// POST /Documents
func (h *DocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var doc domain.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newdoc, err := h.service.Create(context.Background(), &doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newdoc)
}

// GET /Documents
func (h *DocumentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var (
		docType    *string
		status     *string
		employeeID *uuid.UUID
		from, to   *time.Time
	)

	if v := q.Get("type"); v != "" {
		docType = &v
	}
	if v := q.Get("status"); v != "" {
		status = &v
	}
	if v := q.Get("employeeId"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			employeeID = &id
		}
	}
	if v := q.Get("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = &t
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			to = &t
		}
	}

	docs, err := h.service.GetAll(ctx, docType, status, employeeID, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(docs)
}


// GET /Documents/{id}
func (h *DocumentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	doc, err := h.service.GetByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(doc)
}

func (h *DocumentHandler) GetMyDocs(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r.Context())
    role := middlewares.GetUserRole(r.Context())


	docs, err := h.service.GetMyDocs(context.Background(), userID, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(docs)
}

// GET /auth/login
func (h *DocumentHandler) GetMyDocumentsPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../templates/mydocuments.html"))
	tmpl.Execute(w, nil)

}