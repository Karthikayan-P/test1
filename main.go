package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
)

type cust struct {
	CUST_NAME    string    `json:"cust_name"`
	CUST_MAIL    string    `json:"cust_mail"`
	CUST_ADDR    string    `json:"cust_addr"`
	CUST_PHN     int64     `json:"cust_phn"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"password_hash,omitempty"`
	CreatedBy    time.Time `json:"created_by,omitempty"`
}

const (
	internalServerError = "Internal server error"
	passwordLength      = 8
)

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

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
		created_by DATETIME
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	signingKey, err := generateRandomKey()
	if err != nil {
		log.Fatal("Error generating signing key:", err)
	}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerHandler(w, r, db, signingKey)
	})

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, signingKey []byte) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var c cust
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hashedPassword, err := generatePasswordHash(c.Password)
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
	}
	c.PasswordHash = hashedPassword

	currentDate := time.Now().Truncate(24 * time.Hour)
	c.CreatedBy = currentDate

	insertQuery := `INSERT INTO CUSTOMER (cust_name, cust_mail, cust_addr, cust_phn,password_hash,created_by) VALUES (?, ?, ?, ?, ?, ?);`
	result, err := db.Exec(insertQuery, c.CUST_NAME, c.CUST_MAIL, c.CUST_ADDR, c.CUST_PHN, c.PasswordHash, c.CreatedBy)
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

	expirationTime := time.Now().Add(time.Minute * 5)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cust_name": c.CUST_NAME,
		"cust_mail": c.CUST_MAIL,
		"cust_addr": c.CUST_ADDR,
		"cust_phn":  c.CUST_PHN,
		"exp":       expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Registration successful!",
		"token":   tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	decodedToken, err := decodeJWTToken(tokenString, signingKey)
	if err != nil {
		fmt.Println("Error decoding token:", err)
		return
	}

	custName, custMail, custAddr, custPhn, err := getCustomerInfoFromToken(decodedToken)
	if err != nil {
		fmt.Println("Error getting customer information from token:", err)
		return
	}
	fmt.Println("Decoded jwt")
	fmt.Println("Customer Name:", custName)
	fmt.Println("Customer Email:", custMail)
	fmt.Println("Customer Address:", custAddr)
	fmt.Println("Customer Phone Number:", custPhn)

	fmt.Println("token:", tokenString)
}

func decodeJWTToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getCustomerInfoFromToken(token *jwt.Token) (string, string, string, int64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", 0, fmt.Errorf("invalid token claims")
	}

	custName, ok := claims["cust_name"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("invalid customer name in token")
	}

	custMail, ok := claims["cust_mail"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("invalid customer email in token")
	}

	custAddr, ok := claims["cust_addr"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("invalid customer address in token")
	}

	custPhnFloat, ok := claims["cust_phn"].(float64)
	if !ok {
		return "", "", "", 0, fmt.Errorf("invalid customer phone number in token")
	}

	custPhn := int64(custPhnFloat)
	return custName, custMail, custAddr, int64(custPhn), nil
}

func generatePasswordHash(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(hashedPassword), nil
}
