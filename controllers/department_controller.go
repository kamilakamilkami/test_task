package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	_ "project/internal/domain"
	"project/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(s service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: s}
}

// POST /departments
func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string     `json:"name"`
		Code     string     `json:"code"`
		ParentID *uuid.UUID `json:"parentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dept, err := h.service.Create(context.Background(), req.Name, req.Code, req.ParentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(dept)
}

// GET /departments
func (h *DepartmentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	depts, err := h.service.GetAll(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(depts)
}

// GET /departments/{id}
func (h *DepartmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dept, err := h.service.GetByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dept)
}

// PUT /departments/{id}
func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name     string     `json:"name"`
		Code     string     `json:"code"`
		ParentID *uuid.UUID `json:"parentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Update(context.Background(), id, req.Name, req.Code, req.ParentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /departments/{id}
func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
