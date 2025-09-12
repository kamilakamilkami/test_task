package service_domain

import (
	"context"
	"project/internal/domain"
	"project/repository"

	"github.com/google/uuid"
)

type DepartmentService interface {
	Create(ctx context.Context, name string) (*domain.Department, error)
	GetAll(ctx context.Context) ([]domain.Department, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error)
	Update(ctx context.Context, id uuid.UUID, name string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type departmentService struct {
	repo repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

func (s *departmentService) Create(ctx context.Context, name string) (*domain.Department, error) {
	dept := &domain.Department{
		ID:   uuid.New(),
		Name: name,
	}
	err := s.repo.Create(ctx, dept)
	return dept, err
}

func (s *departmentService) GetAll(ctx context.Context) ([]domain.Department, error) {
	return s.repo.GetAll(ctx)
}

func (s *departmentService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *departmentService) Update(ctx context.Context, id uuid.UUID, name string) error {
	dept := &domain.Department{ID: id, Name: name}
	return s.repo.Update(ctx, dept)
}

func (s *departmentService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
