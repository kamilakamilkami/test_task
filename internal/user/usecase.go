package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"project/models"
	"project/utils"
)

type userUseCase struct {
	db *pgxpool.Pool
}

func NewUserUseCase(db *pgxpool.Pool) UseCase {
	return &userUseCase{db: db}
}

func (u *userUseCase) Register(ctx context.Context, user *models.User) error {
	// бизнес-логика: валидация, генерация UUID, хеш пароля
	if user.Email == "" || user.Password == "" {
		return errors.New("email and password required")
	}

	if utils.EmailInvalid(user.Email) {
		return errors.New("invalid email format")
	}

	user.ID = utils.GenerateUUID()
	user.Role = "user"
	user.Password = utils.HashPassword(user.Password)

	// сюда можно передать в репозиторий
	_, err := u.db.Exec(ctx,
		`INSERT INTO users (id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}

func (u *userUseCase) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(ctx,
		`SELECT id, name, email, role FROM users WHERE id=$1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userUseCase) GetAll(ctx context.Context) ([]*models.User, error) {
	rows, err := u.db.Query(ctx, `SELECT id, name, email, role FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *userUseCase) Update(ctx context.Context, user *models.User) error {
	_, err := u.db.Exec(ctx,
		`UPDATE users SET name=$1, email=$2, role=$3 WHERE id=$4`,
		user.Name, user.Email, user.Role, user.ID)
	return err
}

func (u *userUseCase) Delete(ctx context.Context, id string) error {
	_, err := u.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}
