package repository

import (
	"context"
	"project/models"

	_ "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	CheckExists(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email, password, role FROM users WHERE email = $1"
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}

func (r *userRepository) CheckExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}
