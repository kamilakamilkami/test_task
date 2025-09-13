package service

import (
	"context"
	"errors"
	"os"
	"project/dto"
	"project/models"
	"project/repository"
	"project/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (dto.TokenResponse, error)
}

type authUseCase struct {
	repo repository.UserRepository
}

func NewAuthUseCase(r repository.UserRepository) UseCase {
	return &authUseCase{repo: r}
}

func (a *authUseCase) Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error) {
	user, err := a.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.TokenResponse{}, errors.New("invalid email or password")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return dto.TokenResponse{}, errors.New("invalid password")
	}

	tokens, err := utils.GenerateTokens(user.Email, user.Role, os.Getenv("JWT_SECRET"))
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return tokens, nil
}

func (a *authUseCase) Register(ctx context.Context, req dto.RegisterRequest) error {
	exists, _ := a.repo.CheckExists(ctx, req.Email)
	if exists {
		return errors.New("user already exists")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "EMPLOYEE",
	}

	return a.repo.Create(ctx, &user)
}

func (a *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (dto.TokenResponse, error) {
	claims, err := utils.ParseToken(refreshToken, os.Getenv("JWT_SECRET"))
	if err != nil {
		return dto.TokenResponse{}, err
	}

	email := claims["email"].(string)
	role := claims["role"].(string)

	accessToken, accessExpiry, err := utils.GenerateAccessTokenWithExpiry(email, role, os.Getenv("JWT_SECRET"))
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: time.Now().Add(7 * 24 * time.Hour),
	}, nil
}
