package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"project/internal/domain"
	"project/repository"
	"strings"

	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aymerick/raymond"
	"github.com/google/uuid"
)

type DocumentService interface {
	Create(ctx context.Context, doc *domain.Document, userId string) (*domain.Document, error)
	GetAll(ctx context.Context, docType *string, status *string, employeeID *uuid.UUID, from *time.Time, to *time.Time) ([]domain.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error)
	GetMyDocs(ctx context.Context, userId string, role string) ([]domain.Document, error)
	GetFileByID(ctx context.Context, fileID uuid.UUID) (*domain.File, error)
}

type documentService struct {
	repo repository.DocumentRepository
	userRepo repository.UserRepository
	templateRepo repository.TemplateRepository
	sequenceRepo   repository.NumberSequenceRepository
	fileRepo 		repository.FileRepository
}

func NewDocumentService(
    repo repository.DocumentRepository,
    userRepo repository.UserRepository,
    templateRepo repository.TemplateRepository,
	sequenceRepo repository.NumberSequenceRepository,
	fileRepo repository.FileRepository,
) DocumentService {
    return &documentService{
        repo:        repo,
        userRepo:    userRepo,
        templateRepo: templateRepo,
		sequenceRepo: sequenceRepo,
		fileRepo: fileRepo,
    }
}

type Employee struct {
    FullName      string  `json:"fullName"`
    Position      string  `json:"position"`
    HireDate      string  `json:"hireDate"`
    SalaryBase    float64 `json:"salaryBase"`
    SalaryCurrency string `json:"salaryCurrency"`
}

type TemplateData struct {
    Employee      Employee `json:"employee"`
    CompanyName   string   `json:"company_name"`
    IncludeSalary bool     `json:"includeSalary"`
    SalaryInWords string   `json:"salaryInWords"`
    Date          string   `json:"date"`
}


func (s *documentService) Create(ctx context.Context, doc *domain.Document, userId string) (*domain.Document, error) {
	doc.ID = uuid.New()
	employeeID, departmentID, err := s.userRepo.GetEmployeeIdByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee_id for user %s: %w", userId, err)
	}

	empUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id format: %w", err)
	}
	doc.EmployeeID = empUUID

	departmentUUID, err := uuid.Parse(departmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id format: %w", err)
	}
	doc.DepartmentID = departmentUUID

	tempalteID, templateVersion, TemplateBody, err := s.templateRepo.GetByType(ctx, doc.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid template type: %w", err)
	}
	doc.TemplateID = tempalteID
	doc.TemplateVersionID = templateVersion

	number, err := s.sequenceRepo.NextNumber(ctx, doc.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to generate number: %w", err)
	}
	doc.Number = number

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, err
	}

	var templateData TemplateData
	if err := json.Unmarshal([]byte(doc.Data), &templateData); err != nil {
		return nil, fmt.Errorf("invalid document data: %w", err)
	}

	html, err := raymond.Render(TemplateBody, templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	pdfBytes, err := s.FromHtmlToPdf(ctx, html)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	fileID, err := s.SavePdf(ctx, pdfBytes, doc.Number+".pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to save PDF: %w", err)
	}

	doc.FileID = &fileID
	if err := s.repo.UpdateFileID(ctx, doc.ID, fileID); err != nil {
		return nil, fmt.Errorf("failed to update document with fileID: %w", err)
	}

	return doc, nil
}

func (s *documentService) GetAll(ctx context.Context, docType *string, status *string, employeeID *uuid.UUID, from *time.Time, to *time.Time) ([]domain.Document, error) {
	return s.repo.GetAll(ctx, docType, status, employeeID, from, to)
}

func (s *documentService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *documentService) GetMyDocs(ctx context.Context, userId, role string) ([]domain.Document, error) {
	
	if (role == "HR" || role == "ADMIN") {
		return s.GetAll(ctx, nil, nil, nil, nil, nil)
	} else if (role == "MANAGER") {
		return s.repo.GetByDepartmentUserId(ctx, userId)
	} else if (role == "EMPLOYEE") {
		return s.repo.GetByUserId(ctx, userId)
	}
	return []domain.Document{}, fmt.Errorf("There is no such role")

}

func (s *documentService) FromHtmlToPdf(ctx context.Context, html string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF generator: %w", err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(html)))
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	if err := pdfg.Create(); err != nil {
		return nil, fmt.Errorf("failed to create PDF: %w", err)
	}

	return pdfg.Bytes(), nil
}


func (s *documentService) SavePdf(ctx context.Context, pdfBytes []byte, filename string) (uuid.UUID, error) {
	fileID := uuid.New()
	base64Data := base64.StdEncoding.EncodeToString(pdfBytes)

	file := &domain.File{
		ID:       fileID,
		Name:     filename,
		MimeType: "application/pdf",
		Size:     len(pdfBytes),
		Path:     "",
		Base64:   base64Data,
	}

	if err := s.fileRepo.Save(ctx, file); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save PDF: %w", err)
	}

	return fileID, nil
}

func (s *documentService) GetFileByID(ctx context.Context, fileID uuid.UUID) (*domain.File, error) {
    return s.fileRepo.GetByID(ctx, fileID)
}
