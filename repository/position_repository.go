package repository

import (
	"context"
	"project/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PositionRepository interface {
	Create(ctx context.Context, pos *domain.Position) error
	GetAll(ctx context.Context) ([]domain.Position, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Position, error)
	Update(ctx context.Context, pos *domain.Position) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type positionRepository struct {
	db *pgxpool.Pool
}

func NewPositionRepository(db *pgxpool.Pool) PositionRepository {
	return &positionRepository{db: db}
}

func (r *positionRepository) Create(ctx context.Context, pos *domain.Position) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO positions (id, name, code) 
		 VALUES ($1, $2, $3)`,
		pos.ID, pos.Name, pos.Code,
	)
	return err
}

func (r *positionRepository) GetAll(ctx context.Context) ([]domain.Position, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, code FROM positions`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []domain.Position
	for rows.Next() {
		var p domain.Position
		err = rows.Scan(&p.ID, &p.Name, &p.Code)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *positionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Position, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, code FROM positions WHERE id=$1`, id,
	)

	var p domain.Position
	err := row.Scan(&p.ID, &p.Name, &p.Code)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *positionRepository) Update(ctx context.Context, pos *domain.Position) error {
	_, err := r.db.Exec(ctx,
		`UPDATE positions 
		 SET name=$1, code=$2 
		 WHERE id=$3`,
		pos.Name, pos.Code, pos.ID,
	)
	return err
}

func (r *positionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM positions WHERE id=$1`, id)
	return err
}
