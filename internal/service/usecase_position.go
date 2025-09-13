package service

import (
	"context"
	"project/internal/domain"
	"project/repository"

	"github.com/google/uuid"
)

type PositionService interface {
	Create(ctx context.Context, name string, code string) (*domain.Position, error)
	GetAll(ctx context.Context) ([]domain.Position, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Position, error)
	Update(ctx context.Context, id uuid.UUID, name string, code string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type positionService struct {
	repo repository.PositionRepository
}

func NewPositionService(repo repository.PositionRepository) PositionService {
	return &positionService{repo: repo}
}

func (s *positionService) Create(ctx context.Context, name string, code string) (*domain.Position, error) {
	pos := &domain.Position{
		ID:   uuid.New(),
		Name: name,
		Code: code,
	}
	err := s.repo.Create(ctx, pos)
	return pos, err
}

func (s *positionService) GetAll(ctx context.Context) ([]domain.Position, error) {
	return s.repo.GetAll(ctx)
}

func (s *positionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Position, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *positionService) Update(ctx context.Context, id uuid.UUID, name string, code string) error {
	pos := &domain.Position{
		ID:   id,
		Name: name,
		Code: code,
	}
	return s.repo.Update(ctx, pos)
}

func (s *positionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
