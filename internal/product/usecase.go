package product

import (
	"context"
	"errors"
	"project/models"
	"project/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type productUseCase struct {
	db *pgxpool.Pool
}

func NewProductUseCase(db *pgxpool.Pool) UseCase {
	return &productUseCase{db: db}
}

func (u *productUseCase) GetAll(ctx context.Context) ([]*models.Product, error) {
	rows, err := u.db.Query(ctx, `SELECT id, name, price, quantity FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (u *productUseCase) GetByID(ctx context.Context, id string) (*models.Product, error) {
	var p models.Product
	err := u.db.QueryRow(ctx,
		`SELECT id, name, price, quantity FROM products WHERE id=$1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (u *productUseCase) Create(ctx context.Context, p *models.Product) error {
	if p.Name == "" || p.Price <= 0 || p.Quantity < 0 {
		return errors.New("invalid product data")
	}
	p.ID = utils.GenerateUUID()
	_, err := u.db.Exec(ctx,
		`INSERT INTO products (id, name, price, quantity) VALUES ($1, $2, $3, $4)`,
		p.ID, p.Name, p.Price, p.Quantity)
	return err
}

func (u *productUseCase) Update(ctx context.Context, p *models.Product) error {
	_, err := u.db.Exec(ctx,
		`UPDATE products SET name=$1, price=$2, quantity=$3 WHERE id=$4`,
		p.Name, p.Price, p.Quantity, p.ID)
	return err
}

func (u *productUseCase) Delete(ctx context.Context, id string) error {
	_, err := u.db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}
