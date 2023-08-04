package orderFunction

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sample/orderServices"
)

const (
	InternalServerError = "internal server error"
)

// type order orderServices.Order struct {
// 	CustomerID   int      `json:"customer_id"`
// 	ProductIDs   []int    `json:"product_ids"`
// 	ProductNames []string `json:"product_names"`
// }

func RegisterOrderHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Received order registration request.")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var o orderServices.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// orderObj := orderServices.Order{
	// 	CustomerID:   o.CustomerID,
	// 	ProductIDs:   o.ProductIDs,
	// 	ProductNames: o.ProductNames,
	// }

	err = orderServices.InsertOrder(db, o)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		fmt.Println("Database Error:", err)
		return
	}

	response := map[string]string{
		"message": "Order registered successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Println("Order registered successfully!")
	fmt.Printf("CustomerID: %d\n", o.CustomerID)
	fmt.Printf("ProductIDs: %v\n", o.ProductIDs)
	fmt.Printf("ProductNames: %v\n", o.ProductNames)
}
