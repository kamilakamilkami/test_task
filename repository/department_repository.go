package repository

import (
	"context"
	"project/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DepartmentRepository interface {
	Create(ctx context.Context, dept *domain.Department) error
	GetAll(ctx context.Context) ([]domain.Department, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error)
	Update(ctx context.Context, dept *domain.Department) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type departmentRepository struct {
	db *pgxpool.Pool
}

func NewDepartmentRepository(db *pgxpool.Pool) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, dept *domain.Department) error {
	_, err := r.db.Exec(ctx, "INSERT INTO departments (id, name) VALUES ($1, $2)", dept.ID, dept.Name)
	return err
}

func (r *departmentRepository) GetAll(ctx context.Context) ([]domain.Department, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name FROM departments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []domain.Department
	for rows.Next() {
		var d domain.Department
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, nil
}

func (r *departmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	row := r.db.QueryRow(ctx, "SELECT id, name FROM departments WHERE id=$1", id)

	var d domain.Department
	if err := row.Scan(&d.ID, &d.Name); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *departmentRepository) Update(ctx context.Context, dept *domain.Department) error {
	_, err := r.db.Exec(ctx, "UPDATE departments SET name=$1 WHERE id=$2", dept.Name, dept.ID)
	return err
}

func (r *departmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM departments WHERE id=$1", id)
	return err
}
