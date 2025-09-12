package service

import (
	"context"
	"project/internal/domain"
	"project/repository"

	"github.com/google/uuid"
	"time"
)

type DocumentService interface {
	Create(ctx context.Context, doc *domain.Document) (*domain.Document, error)
	GetAll(ctx context.Context, docType *string, status *string, employeeID *uuid.UUID, from *time.Time, to *time.Time) ([]domain.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error)
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
