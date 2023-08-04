package ordergetFunction

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sample/ordergetService"
	"strconv"
)

const (
	InternalServerError = "internal server error"
)

func GetOrderDetailsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Received get order details request.")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customerIDStr := r.URL.Query().Get("customer_id")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	orderDetails, err := ordergetService.GetOrderDetails(db,customerID)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		fmt.Println("Database Error:", err)
		return
	}

	fmt.Println("List of details for order")
	for _, orderDetail := range orderDetails {
		fmt.Printf("OrderID: %d\n", orderDetail.OrderID)
		fmt.Printf("CustomerID: %d\n", orderDetail.CustomerID)
		fmt.Printf("CustomerName: %s\n", orderDetail.CustomerName)
		fmt.Printf("ProductIDs: %v\n", orderDetail.ProductIDs)
		fmt.Printf("ProductNames: %v\n", orderDetail.ProductNames)
		fmt.Println()
	}

	response := map[string]string{
		"message": "Order details fetched successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
