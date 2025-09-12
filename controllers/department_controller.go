package controllers

import (
	"encoding/json"
	"net/http"
	service "project/internal/service_domain"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(service service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dept, err := h.service.Create(r.Context(), input.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	depts, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(depts)
}

func (h *DepartmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	dept, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.Update(r.Context(), id, input.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
