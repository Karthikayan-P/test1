package productServices

import (
	"database/sql"
)

func CreateProductTable(db *sql.DB) error {
	createProductTable := `CREATE TABLE IF NOT EXISTS PRODUCT(
		product_id INT AUTO_INCREMENT PRIMARY KEY,
		prod_name VARCHAR(100),
		quantity INT,
		availability VARCHAR(3),
		product_amount FLOAT,
		customer_id INT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (customer_id) REFERENCES CUSTOMER(cust_id)
	);`

	_, err := db.Exec(createProductTable)
	if err != nil {
		return err
	}

	return nil
}
