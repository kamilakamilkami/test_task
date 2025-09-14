package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"project/internal/domain"
	"project/internal/service"

	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"html/template"
	"log"
	"project/middlewares"
)

type DocumentHandler struct {
	service service.DocumentService
}

func NewDocumentHandler(s service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: s}
}

// POST /Documents
func (h *DocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
    userID := middlewares.GetUserID(r.Context())

    var doc domain.Document
    if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
        log.Printf("❌ Decode error: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    newdoc, err := h.service.Create(r.Context(), &doc, userID) 
    if err != nil {
        log.Printf("❌ Create error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(newdoc); err != nil {
        log.Printf("❌ Encode response error: %v", err)
    }
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

// GET /mydocuments
func (h *DocumentHandler) GetMyDocumentsPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/mydocuments.html"))
	tmpl.Execute(w, nil)

}

// GET /createcertificate
func (h *DocumentHandler) GetCreateCertificatePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/createcertificate.html"))
	tmpl.Execute(w, nil)

}
type PreviewData struct {
	ID            string
	FullName      string
	Position      string
	Department    string
	EmployedAt    string
	Salary        *float64
	IncludeSalary bool
	ExpiresInDays int
}

// GET /previewcertificate
func (h *DocumentHandler) PreviewCertificatePage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	doc, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	var docData map[string]interface{}
	if err := json.Unmarshal([]byte(doc.Data), &docData); err != nil {
		http.Error(w, "Invalid document data", http.StatusInternalServerError)
		return
	}

	var docMeta map[string]interface{}
	if err := json.Unmarshal([]byte(doc.Meta), &docMeta); err != nil {
		http.Error(w, "Invalid document meta", http.StatusInternalServerError)
		return
	}

	employeeMap := make(map[string]interface{})
	if e, ok := docData["employee"].(map[string]interface{}); ok {
		employeeMap = e
	}

	fullName, _ := employeeMap["fullName"].(string)
	position, _ := employeeMap["position"].(string)
	employedAt, _ := employeeMap["hireDate"].(string)
	var salary *float64
	includeSalary := false
	if val, ok := employeeMap["salaryBase"].(float64); ok {
		salary = &val
		includeSalary = true
	}

	department, _ := docData["department"].(string)

	preview := PreviewData{
		ID:            doc.ID.String(),
		FullName:      fullName,
		Position:      position,
		Department:    department,
		EmployedAt:    employedAt,
		Salary:        salary,
		IncludeSalary: includeSalary,
	}

	if expires, ok := docMeta["expiresInDays"].(float64); ok {
		preview.ExpiresInDays = int(expires)
	}

	tmpl := template.Must(template.ParseFiles("templates/preview.html"))
	if err := tmpl.Execute(w, preview); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (h *DocumentHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    docID, err := uuid.Parse(params["id"])
    if err != nil {
        http.Error(w, "Invalid document ID", http.StatusBadRequest)
        return
    }

    // Получаем документ
    doc, err := h.service.GetByID(r.Context(), docID)
    if err != nil {
        http.Error(w, "Document not found", http.StatusNotFound)
        return
    }

    if doc.FileID == nil {
        http.Error(w, "PDF not generated yet", http.StatusNotFound)
        return
    }

    // Получаем файл из БД
    file, err := h.service.GetFileByID(r.Context(), *doc.FileID)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Декодируем base64 в байты
    pdfBytes, err := base64.StdEncoding.DecodeString(file.Base64)
    if err != nil {
        http.Error(w, "Failed to decode PDF", http.StatusInternalServerError)
        return
    }

    // Отдаём пользователю
    w.Header().Set("Content-Type", file.MimeType)
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", file.Name))
    w.Write(pdfBytes)
}



