package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sample/customerFunction"
	"sample/customerServices"
	"sample/orderFunction"
	"sample/orderServices"
	"sample/ordergetFunction"
	"sample/productFunction"
	"sample/productServices"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:3307)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable := `CREATE TABLE IF NOT EXISTS CUSTOMER(
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

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	err = customerServices.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}

	err = productServices.CreateProductTable(db)
	if err != nil {
		log.Fatal(err)
	}

	err = orderServices.CreateOrderTable(db)
	if err != nil {
		log.Fatal(err)
	}

	signingKey, err := customerFunction.GenerateRandomKey()
	if err != nil {
		log.Fatal("Error generating signing key:", err)
	}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		customerFunction.RegisterHandler(w, r, db, signingKey)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		customerFunction.LoginHandler(w, r, db, signingKey)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		customerFunction.DeleteHandler(w, r, db)
	})

	http.HandleFunc("/cd", func(w http.ResponseWriter, r *http.Request) {
		productFunction.RegisterProductHandler(w, r, db)
	})

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		orderFunction.RegisterOrderHandler(w, r, db)
	})

	http.HandleFunc("/get-order-details", func(w http.ResponseWriter, r *http.Request) {
		ordergetFunction.GetOrderDetailsHandler(w, r, db)
	})

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
