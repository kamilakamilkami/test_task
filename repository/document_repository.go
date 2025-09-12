package repository

import (
	"context"
	"project/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"fmt"
	"time"
)

type DocumentRepository interface {
	Create(ctx context.Context, emp *domain.Document) error
	GetAll(ctx context.Context, docType *string, status *string, employeeID *uuid.UUID, from *time.Time, to *time.Time) ([]domain.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error)
}

type documentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, doc *domain.Document) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO documents 
		(id, type, employee_id, template_id, template_version, number, date, status, file_id, data, meta) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, $11)`,
		doc.ID,
		doc.Type,
		doc.EmployeeID,
		doc.TemplateID,
		doc.TemplateVersionID,
		doc.Number,
		doc.Date,
		doc.Status,
		doc.FileID,
		doc.Data, 
		doc.Meta, 
	)
	return err
}


func (r *documentRepository) GetAll(
    ctx context.Context,
    docType *string,
    status *string,
    employeeID *uuid.UUID,
    from *time.Time,
    to *time.Time,
) ([]domain.Document, error) {
    // Базовый запрос
    query := `
        SELECT id, type, employee_id, template_id, template_version, number, date, status, file_id, data, meta
        FROM documents
        WHERE 1=1
    `
    args := []interface{}{}
    argIdx := 1

    if docType != nil {
        query += fmt.Sprintf(" AND type = $%d", argIdx)
        args = append(args, *docType)
        argIdx++
    }

    if status != nil {
        query += fmt.Sprintf(" AND status = $%d", argIdx)
        args = append(args, *status)
        argIdx++
    }

    if employeeID != nil {
        query += fmt.Sprintf(" AND employee_id = $%d", argIdx)
        args = append(args, *employeeID)
        argIdx++
    }

    if from != nil {
        query += fmt.Sprintf(" AND date >= $%d", argIdx)
        args = append(args, *from)
        argIdx++
    }

    if to != nil {
        query += fmt.Sprintf(" AND date <= $%d", argIdx)
        args = append(args, *to)
        argIdx++
    }

    query += " ORDER BY date DESC"

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var docs []domain.Document
    for rows.Next() {
        var doc domain.Document
        err := rows.Scan(
            &doc.ID,
            &doc.Type,
            &doc.EmployeeID,
			&doc.TemplateID,
            &doc.TemplateVersionID,
            &doc.Number,
            &doc.Date,
            &doc.Status,
            &doc.FileID,
            &doc.Data,
            &doc.Meta,
        )
        if err != nil {
            return nil, err
        }
        docs = append(docs, doc)
    }

    return docs, nil
}



func (r *documentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error) {
	var doc domain.Document
	err := r.db.QueryRow(ctx, `
		SELECT id, type, employee_id, template_id, template_version, number, date, status, file_id, data, meta
		FROM documents 
		WHERE id = $1
	`, id).Scan(
		&doc.ID,
		&doc.Type,
		&doc.EmployeeID,
		&doc.TemplateID,
		&doc.TemplateVersionID,
		&doc.Number,
		&doc.Date,
		&doc.Status,
		&doc.FileID,
		&doc.Data,
		&doc.Meta,
	)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}


