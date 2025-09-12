package auth

import (
	"context"
	"errors"
	"os"
	"project/dto"
	"project/models"
	"project/repository"
	"project/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	db *pgxpool.Pool
}

func NewAuthUseCase(db *pgxpool.Pool) UseCase {
	return &authUseCase{db: db}
}

func (a *authUseCase) Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error) {
	user := repository.GetUser(ctx, a.db, req.Email)
	if user == nil {
		return dto.TokenResponse{}, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.TokenResponse{}, errors.New("invalid password")
	}

	tokens, err := utils.GenerateTokens(user.Email, user.Role, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return tokens, nil
}

func (a *authUseCase) Register(ctx context.Context, user *models.User) error {
	if repository.CheckUserExists(ctx, a.db, user.Email) {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.ID = utils.GenerateUUID()
	user.Password = string(hashedPassword)
	user.Role = "user"

	return repository.CreateUser(ctx, a.db, user)
}

func (a *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (dto.TokenResponse, error) {
	claims, err := utils.ParseToken(refreshToken, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return dto.TokenResponse{}, err
	}

	email := claims["email"].(string)
	role := claims["role"].(string)

	accessToken, accessExpiry, err := utils.GenerateAccessTokenWithExpiry(email, role, os.Getenv("JWT_SECRET_KEY"))
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
