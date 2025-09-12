package product

import (
	"context"
	"project/models"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]*models.Product, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id string) error
}
