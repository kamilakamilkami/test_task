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
	_, err := r.db.Exec(ctx,
		`INSERT INTO departments (id, name, code, parent_id) 
		 VALUES ($1, $2, $3, $4)`,
		dept.ID, dept.Name, dept.Code, dept.ParentID,
	)
	return err
}

func (r *departmentRepository) GetAll(ctx context.Context) ([]domain.Department, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, code, parent_id FROM departments`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var depts []domain.Department
	for rows.Next() {
		var d domain.Department
		err = rows.Scan(&d.ID, &d.Name, &d.Code, &d.ParentID)
		if err != nil {
			return nil, err
		}
		depts = append(depts, d)
	}
	return depts, nil
}

func (r *departmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, code, parent_id FROM departments WHERE id=$1`, id,
	)

	var d domain.Department
	err := row.Scan(&d.ID, &d.Name, &d.Code, &d.ParentID)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *departmentRepository) Update(ctx context.Context, dept *domain.Department) error {
	_, err := r.db.Exec(ctx,
		`UPDATE departments 
		 SET name=$1, code=$2, parent_id=$3 
		 WHERE id=$4`,
		dept.Name, dept.Code, dept.ParentID, dept.ID,
	)
	return err
}

func (r *departmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM departments WHERE id=$1`, id)
	return err
}
