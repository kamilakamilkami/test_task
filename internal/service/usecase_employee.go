package service

import (
	"context"
	"project/internal/domain"
	"project/repository"

	"github.com/google/uuid"
)

type EmployeeService interface {
	Create(ctx context.Context, emp *domain.Employee) (*domain.Employee, error)
	GetAll(ctx context.Context) ([]domain.Employee, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	Update(ctx context.Context, emp *domain.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) Create(ctx context.Context, emp *domain.Employee) (*domain.Employee, error) {
	emp.ID = uuid.New()
	if err := s.repo.Create(ctx, emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *employeeService) GetAll(ctx context.Context) ([]domain.Employee, error) {
	return s.repo.GetAll(ctx)
}

func (s *employeeService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *employeeService) Update(ctx context.Context, emp *domain.Employee) error {
	return s.repo.Update(ctx, emp)
}

func (s *employeeService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
