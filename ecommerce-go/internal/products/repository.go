package products

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Stock       int
}

func CreateTable(conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC(10,2) NOT NULL,
		stock INT NOT NULL
	);`
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error creando tabla products: %v", err)
	}
	fmt.Println("Tabla 'products' lista")
	return nil
}

func GetAll(conn *pgx.Conn) ([]Product, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, name, description, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func Create(conn *pgx.Conn, p Product) error {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4)",
		p.Name, p.Description, p.Price, p.Stock)
	return err
}
