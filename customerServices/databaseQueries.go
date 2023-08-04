package customerServices

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func CreateTables(db *sql.DB) error {
	createCustomerTable := `CREATE TABLE IF NOT EXISTS CUSTOMER(
		cust_id INT AUTO_INCREMENT PRIMARY KEY,
		cust_name VARCHAR(50),
		cust_mail VARCHAR(100),
		cust_addr VARCHAR(200), 
		cust_phn BIGINT,
		password_hash VARCHAR(255),
		created_by DATETIME,
		active INT DEFAULT 1,
		update_at DATETIME
	);`

	_, err := db.Exec(createCustomerTable)
	if err != nil {
		return err
	}

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

	_, err = db.Exec(createProductTable)
	if err != nil {
		return err
	}

	return nil
}
