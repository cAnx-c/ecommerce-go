package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cAnx-c/ecommerce-go/internal/database"
	"github.com/cAnx-c/ecommerce-go/internal/products"
)

func main() {

	conn := database.Connect()
	defer conn.Close(context.Background())

	if err := products.CreateTable(conn); err != nil {
		log.Fatalf("Error creando tabla: %v", err)
	}

	products.RegisterHandlers()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenido al Sistema de Gestión e-commerce en Go")
	})

	fmt.Println("Servidor e-commerce iniciado en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
