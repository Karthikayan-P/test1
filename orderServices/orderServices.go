package orderServices

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Order struct {
	OrderID      int      `json:"order_id,omitempty"`
	CustomerID   int      `json:"customer_id"`
	ProductIDs   []int    `json:"product_ids"`
	ProductNames []string `json:"product_names"`
}

func CreateOrderTable(db *sql.DB) error {
	createOrderTable := `CREATE TABLE IF NOT EXISTS ORDER_TABLE(
		order_id INT AUTO_INCREMENT PRIMARY KEY,
		customer_id INT,
		product_ids JSON,
		product_names JSON,
		order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (customer_id) REFERENCES CUSTOMER(cust_id)
	);`

	_, err := db.Exec(createOrderTable)
	if err != nil {
		return err
	}

	return nil
}

func InsertOrder(db *sql.DB, o Order) error {
	productIDsJSON, err := json.Marshal(o.ProductIDs)
	if err != nil {
		return err
	}

	productNamesJSON, err := json.Marshal(o.ProductNames)
	if err != nil {
		return err
	}

	insertQuery := `INSERT INTO ORDER_TABLE (customer_id, product_ids, product_names) VALUES (?, ?, ?);`
	result, err := db.Exec(insertQuery, o.CustomerID, productIDsJSON, productNamesJSON)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected")
	}

	return nil
}
