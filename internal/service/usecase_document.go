package service

import (
	"context"
	"fmt"
	"project/internal/domain"
	"project/repository"

	"time"

	"github.com/google/uuid"
)

type DocumentService interface {
	Create(ctx context.Context, doc *domain.Document, userId string) (*domain.Document, error)
	GetAll(ctx context.Context, docType *string, status *string, employeeID *uuid.UUID, from *time.Time, to *time.Time) ([]domain.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error)
	GetMyDocs(ctx context.Context, userId string, role string) ([]domain.Document, error)
}

type documentService struct {
	repo repository.DocumentRepository
	userRepo repository.UserRepository
	templateRepo repository.TemplateRepository
	sequenceRepo   repository.NumberSequenceRepository
}

func NewDocumentService(
    repo repository.DocumentRepository,
    userRepo repository.UserRepository,
    templateRepo repository.TemplateRepository,
	sequenceRepo repository.NumberSequenceRepository,
) DocumentService {
    return &documentService{
        repo:        repo,
        userRepo:    userRepo,
        templateRepo: templateRepo,
		sequenceRepo: sequenceRepo,
    }
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

	tempalteID, templateVersion, err := s.templateRepo.GetByType(ctx, doc.Type)
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