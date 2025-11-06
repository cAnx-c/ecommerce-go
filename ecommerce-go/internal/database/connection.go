package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Connect() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:admin@localhost:5432/ecommerce_db")
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	fmt.Println("Conectado exitosamente a PostgreSQL con pgx")
	return conn
}
