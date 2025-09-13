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
	Create(ctx context.Context, doc *domain.Document) (*domain.Document, error)
	GetAll(ctx context.Context, docType *string, status *string, employeeID *uuid.UUID, from *time.Time, to *time.Time) ([]domain.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error)
	GetMyDocs(ctx context.Context, userId string, role string) ([]domain.Document, error)
}

type documentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) DocumentService {
	return &documentService{repo: repo}
}

func (s *documentService) Create(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
	// doc.ID = uuid.New()
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