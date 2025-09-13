package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
	"fmt"
)

type NumberSequenceRepository interface {
    NextNumber(ctx context.Context, docType string) (string, error)
}

type numberSequenceRepository struct {
    db *pgxpool.Pool
}

func NewNumberSequenceRepository(db *pgxpool.Pool) NumberSequenceRepository {
    return &numberSequenceRepository{db: db}
}

func (r *numberSequenceRepository) NextNumber(ctx context.Context, docType string) (string, error) {
    tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return "", err
    }
    defer tx.Rollback(ctx)

    var lastValue int
    err = tx.QueryRow(ctx, `
        SELECT last_value 
        FROM number_sequences 
        WHERE doc_type = $1 FOR UPDATE
    `, docType).Scan(&lastValue)

    if err == pgx.ErrNoRows {
        // 👇 создаём сразу с last_value = 1
        _, err = tx.Exec(ctx, `
            INSERT INTO number_sequences (id, doc_type, pattern, last_value)
            VALUES ($1, $2, $3, 1)
        `, uuid.New(), docType, "CERT-{yyyyMMdd}-{seq}")
        if err != nil {
            return "", err
        }
        lastValue = 1
    } else if err != nil {
        return "", err
    } else {
        lastValue++ // 👈 увеличиваем, а не вставляем заново
        _, err = tx.Exec(ctx, `
            UPDATE number_sequences 
            SET last_value = $1, updated_at = now()
            WHERE doc_type = $2
        `, lastValue, docType)
        if err != nil {
            return "", err
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return "", err
    }

    today := time.Now().Format("20060102")
    return fmt.Sprintf("CERT-%s-%03d", today, lastValue), nil
}
