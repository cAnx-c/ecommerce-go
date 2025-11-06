package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cAnx-c/ecommerce-go/internal/database"
)

func RegisterHandlers() {
	http.HandleFunc("/products", productsHandler)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	case http.MethodPut:
		updateProduct(w, r)
	case http.MethodDelete:
		deleteProduct(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	conn := database.Connect()
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT id, name, description, price, stock FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock)
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	conn := database.Connect()
	defer conn.Close(context.Background())

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	_, err := conn.Exec(context.Background(),
		"INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4)",
		p.Name, p.Description, p.Price, p.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Producto creado correctamente")
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	conn := database.Connect()
	defer conn.Close(context.Background())

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Falta el parámetro ?id=", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	_, err = conn.Exec(context.Background(),
		"UPDATE products SET name=$1, description=$2, price=$3, stock=$4 WHERE id=$5",
		p.Name, p.Description, p.Price, p.Stock, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Producto actualizado correctamente")
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	conn := database.Connect()
	defer conn.Close(context.Background())

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		parts := strings.Split(r.URL.Path, "/")
		idStr = parts[len(parts)-1]
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = conn.Exec(context.Background(), "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Producto eliminado correctamente")
}
