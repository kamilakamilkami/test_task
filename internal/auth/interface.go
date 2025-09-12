package auth

import (
	"context"
	"project/dto"
	"project/models"
)

type UseCase interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error)
	Register(ctx context.Context, user *models.User) error
	RefreshToken(ctx context.Context, refreshToken string) (dto.TokenResponse, error)
}
