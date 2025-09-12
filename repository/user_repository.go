package repository

import (
	"context"
	"project/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserByID(ctx context.Context, db *pgxpool.Pool, id string) *models.User {
	var user models.User
	query := "SELECT id, name, email, role FROM users WHERE id = $1"
	err := db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return nil
	}
	return &user
}

func GetUsers(ctx context.Context, db *pgxpool.Pool) []*models.User {
	rows, err := db.Query(ctx, "SELECT id, name, email, role FROM users")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			return nil
		}
		users = append(users, &user)
	}
	return users
}

func UpdateUser(ctx context.Context, db *pgxpool.Pool, user *models.User) error {
	query := "UPDATE users SET name = $1, email = $2, role = $3 WHERE id = $4"
	_, err := db.Exec(ctx, query, user.Name, user.Email, user.Role, user.ID)
	return err
}

func DeleteUser(ctx context.Context, db *pgxpool.Pool, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(ctx, query, id)
	return err
}
