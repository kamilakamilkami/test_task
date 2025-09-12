package repository

import (
	"context"
	"project/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetProducts(ctx context.Context, db *pgxpool.Pool) []*models.Product {
	rows, err := db.Query(ctx, "SELECT id, name, price, quantity FROM products")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
			return nil
		}
		products = append(products, &product)
	}
	return products
}

func GetProductByID(ctx context.Context, db *pgxpool.Pool, id string) *models.Product {
	var product models.Product
	err := db.QueryRow(ctx, "SELECT id, name, price, quantity FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return nil
	}
	return &product
}

func CreateProduct(ctx context.Context, db *pgxpool.Pool, product *models.Product) error {
	query := "INSERT INTO products (id, name, price, quantity) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(ctx, query, product.ID, product.Name, product.Price, product.Quantity)
	return err
}

func UpdateProduct(ctx context.Context, db *pgxpool.Pool, product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, quantity = $3 WHERE id = $4"
	_, err := db.Exec(ctx, query, product.Name, product.Price, product.Quantity, product.ID)
	return err
}

func DeleteProduct(ctx context.Context, db *pgxpool.Pool, id string) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := db.Exec(ctx, query, id)
	return err
}
