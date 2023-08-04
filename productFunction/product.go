package productFunction

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type product struct {
	ProductID     int     `json:"product_id,omitempty"`
	ProdName      string  `json:"prod_name"`
	Quantity      int     `json:"quantity"`
	Availability  string  `json:"availability"`
	ProductAmount float64 `json:"product_amount"`
	CustomerID    int     `json:"customer_id,omitempty"`
}

const (
	internalServerError = "Internal Server Error"
)

func RegisterProductHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Received product registration request.")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	insertQuery := `INSERT INTO PRODUCT (prod_name, quantity, availability, product_amount, customer_id) VALUES (?, ?, ?, ?, ?);`
	result, err := db.Exec(insertQuery, p.ProdName, p.Quantity, p.Availability, p.ProductAmount, p.CustomerID)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		fmt.Println("Database Error:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		fmt.Println("RowsAffected Error:", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		fmt.Println("No rows were affected")
		return
	}

	response := map[string]string{
		"message": "Product registered successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
