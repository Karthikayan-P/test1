package ordergetService

import (
	"database/sql"
	"encoding/json"
)

type OrderDetail struct {
	OrderID      int      `json:"order_id"`
	CustomerID   int      `json:"customer_id"`
	CustomerName string   `json:"customer_name"`
	ProductIDs   []int    `json:"product_ids"`
	ProductNames []string `json:"product_names"`
}

func GetOrderDetails(db *sql.DB, customerID int) ([]OrderDetail, error) {
	query := `
		SELECT 
			o.order_id,
			o.customer_id,
			c.cust_name,
			o.product_ids,
			o.product_names
		FROM
			ORDER_TABLE o
		INNER JOIN
			CUSTOMER c ON o.customer_id = c.cust_id
		WHERE
			o.customer_id= ?;
	`

	rows, err := db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderDetails []OrderDetail
	for rows.Next() {
		var orderDetail OrderDetail
		var productIDsJSON, productNamesJSON []byte

		err := rows.Scan(
			&orderDetail.OrderID,
			&orderDetail.CustomerID,
			&orderDetail.CustomerName,
			&productIDsJSON,
			&productNamesJSON,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(productIDsJSON, &orderDetail.ProductIDs)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(productNamesJSON, &orderDetail.ProductNames)
		if err != nil {
			return nil, err
		}

		orderDetails = append(orderDetails, orderDetail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderDetails, nil
}
