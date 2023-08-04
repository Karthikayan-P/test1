package customerFunction

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
	Active       int       `json:"active,omitempty"`
	updateat     time.Time `json:"update_at,omitempty"`
}

const (
	internalServerError = "Internal server error"
	passwordLength      = 8
)

func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, signingKey []byte) {
	log.Println("Received registration request.")
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

	hashedPassword, err := GeneratePasswordHash(c.Password)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
	c.PasswordHash = hashedPassword

	currentDate := time.Now().Truncate(24 * time.Hour)
	c.CreatedBy = currentDate

	c.Active = 1

	c.updateat = time.Now()

	insertQuery := `INSERT INTO CUSTOMER (cust_name, cust_mail, cust_addr, cust_phn,password_hash,created_by,active,update_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	result, err := db.Exec(insertQuery, c.CUST_NAME, c.CUST_MAIL, c.CUST_ADDR, c.CUST_PHN, c.PasswordHash, c.CreatedBy, c.Active, c.updateat)
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

	decodedToken, err := DecodeJWTToken(tokenString, signingKey)
	if err != nil {
		fmt.Println("Error decoding token:", err)
		return
	}

	custName, custMail, custAddr, custPhn, err := GetCustomerInfoFromToken(decodedToken)
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

func DecodeJWTToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
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

func GetCustomerInfoFromToken(token *jwt.Token) (string, string, string, int64, error) {
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

func GeneratePasswordHash(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
func ValidatePasswordHash(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, signingKey []byte) {
	log.Println("Received login request.")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData struct {
		CUST_MAIL string `json:"cust_mail"`
		Password  string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err = db.QueryRow("SELECT password_hash FROM CUSTOMER WHERE cust_mail = ?", loginData.CUST_MAIL).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = ValidatePasswordHash(loginData.Password, hashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginData)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deleteData struct {
		CUST_MAIL string `json:"cust_mail"`
	}

	err := json.NewDecoder(r.Body).Decode(&deleteData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedAt := time.Now()
	updateQuery := `UPDATE CUSTOMER SET active = 0, update_at = ? WHERE cust_mail = ?;`
	result, err := db.Exec(updateQuery, updatedAt, deleteData.CUST_MAIL)
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
		http.Error(w, "Customer not found", http.StatusNotFound)
		fmt.Println("No rows were affected")
		return
	}

	response := map[string]string{
		"message": "Customer deleted successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
