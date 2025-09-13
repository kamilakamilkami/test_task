package repository

import (
	"context"
	"encoding/base64"
	"project/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository interface {
	Save(ctx context.Context, file *domain.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error)
}

type fileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Save(ctx context.Context, file *domain.File) error {
	query := `
		INSERT INTO files (id, name, mime_type, size, path, base64, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
	`
	_, err := r.db.Exec(ctx, query, file.ID, file.Name, file.MimeType, file.Size, file.Path, file.Base64)
	return err
}

func (r *fileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	var f domain.File
	query := `SELECT id, name, mime_type, size, path, base64, created_at FROM files WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.Name, &f.MimeType, &f.Size, &f.Path, &f.Base64, &f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *fileRepository) GetByIDBase(ctx context.Context, id uuid.UUID) (*domain.File, error) {
    var f domain.File
    query := `SELECT id, name, mime_type, size, path, base64, created_at FROM files WHERE id = $1`
    err := r.db.QueryRow(ctx, query, id).Scan(
        &f.ID, &f.Name, &f.MimeType, &f.Size, &f.Path, &f.Base64, &f.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    fBytes, _ := base64.StdEncoding.DecodeString(f.Base64)
    f.Base64 = string(fBytes)
    return &f, nil
}
