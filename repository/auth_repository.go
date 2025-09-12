package repository

import (
	"context"
	"project/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckUserExists(ctx context.Context, db *pgxpool.Pool, email string) bool {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return err == nil && exists
}

func GetUser(ctx context.Context, db *pgxpool.Pool, email string) *models.User {
	var user models.User
	query := "SELECT id, email, password, role FROM users WHERE email = $1"
	err := db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil
	}
	return &user
}

func CreateUser(ctx context.Context, db *pgxpool.Pool, user *models.User) error {
	query := "INSERT INTO users (id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}
