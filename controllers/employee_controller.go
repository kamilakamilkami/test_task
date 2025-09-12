package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"project/internal/domain"
	"project/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(s service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

// POST /employees
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var emp domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newEmp, err := h.service.Create(context.Background(), &emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newEmp)
}

// GET /employees
func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	emps, err := h.service.GetAll(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(emps)
}

// GET /employees/{id}
func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	emp, err := h.service.GetByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(emp)
}

// PUT /employees/{id}
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var emp domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	emp.ID = id

	if err := h.service.Update(context.Background(), &emp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /employees/{id}
func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
