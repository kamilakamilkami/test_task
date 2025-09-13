package repository

import (
	"context"
	"project/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"fmt"
	"time"
	"strings"
)

type TemplateRepository interface {
	Create(ctx context.Context, template *domain.Template) error
	GetAll(ctx context.Context, docType *string, active *bool) ([]domain.Template, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error)
	Update(ctx context.Context, id uuid.UUID, template *domain.Template) (*domain.Template, error)
}

type templateRepository struct {
	db *pgxpool.Pool
}

func NewTemplateRepository(db *pgxpool.Pool) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(ctx context.Context, t *domain.Template) error {
    _, err := r.db.Exec(ctx,
        `INSERT INTO templates 
        (name, type, engine, body, placeholders, version, is_active, created_by, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
        t.Name,
        t.Type,
        t.Engine,
        t.Body,
        t.Placeholders, 
        t.Version,
        t.IsActive,
        t.CreatedBy,
        time.Now(),     
    )
    return err
}


func (r *templateRepository) GetAll(ctx context.Context, templateType *string, active *bool) ([]domain.Template, error) {
    query := `SELECT id, name, type, engine, body, placeholders, version, is_active, created_by, created_at FROM templates`
    args := []interface{}{}
    conditions := []string{}

    if templateType != nil {
        conditions = append(conditions, "type = $"+fmt.Sprint(len(args)+1))
        args = append(args, *templateType)
    }
    if active != nil {
        conditions = append(conditions, "is_active = $"+fmt.Sprint(len(args)+1))
        args = append(args, *active)
    }

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var templates []domain.Template
    for rows.Next() {
        var t domain.Template
        err := rows.Scan(
            &t.ID, &t.Name, &t.Type, &t.Engine, &t.Body, &t.Placeholders,
            &t.Version, &t.IsActive, &t.CreatedBy, &t.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        templates = append(templates, t)
    }
    return templates, nil
}




func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error) {
	var template domain.Template
	err := r.db.QueryRow(ctx, `
		SELECT id, name, type, engine, body, placeholders, version, is_active, created_by, created_at
		FROM templates 
		WHERE id = $1
	`, id).Scan(
		&template.ID,
		&template.Name,
        &template.Type,
        &template.Engine,
        &template.Body,
        &template.Placeholders, 
        &template.Version,
        &template.IsActive,
        &template.CreatedBy,
        &template.CreatedAt, 
	)
	if err != nil {
		return nil, err
	}
	return &template, nil
}


func (r *templateRepository) Update(ctx context.Context, id uuid.UUID, t *domain.Template) (*domain.Template, error) {

	_, err := r.db.Exec(ctx, `
        UPDATE templates
        SET name = $1,
            type = $2,
            engine = $3,
            body = $4,
            placeholders = $5,
            version = version + 1,
            is_active = $6
        WHERE id = $7
        RETURNING id, name, type, engine, body, placeholders, version, is_active, created_by, created_at
    `,
        t.Name, t.Type, t.Engine, t.Body, t.Placeholders, t.IsActive, id,
    )
    if err != nil {
        return nil, err
    }

    var updated domain.Template
    err = r.db.QueryRow(ctx, `
        SELECT id, name, type, engine, body, placeholders, version, is_active, created_by, created_at
        FROM templates WHERE id = $1
    `, id).Scan(
        &updated.ID, &updated.Name, &updated.Type, &updated.Engine, &updated.Body, &updated.Placeholders,
        &updated.Version, &updated.IsActive, &updated.CreatedBy, &updated.CreatedAt,
    )
    if err != nil {
        return nil, err
    }

    return &updated, nil

}