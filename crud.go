package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type employee struct {
	ID     int
	Name   string
	Age    int
	Salary int
}

func main() {
	db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:3307)/classicmodels")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable := `CREATE TABLE IF NOT EXISTS employe(
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50),
		age INT,
		salary INT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	employees := []employee{
		{Name: "Karthik", Age: 25, Salary: 30000},
		{Name: "Shankari", Age: 24, Salary: 49000},
	}
	for _, emp := range employees {
		err = createEmployee(db, emp)
		if err != nil {
			log.Fatal(err)
		}
	}

	employees, err = getEmployees(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Employees:")
	for _, emp := range employees {
		fmt.Printf("ID: %d, Name: %s, Age: %d, Salary: %d\n", emp.ID, emp.Name, emp.Age, emp.Salary)
	}

	employeeID := 1
	newAge := 31
	err = updateEmployeeAge(db, employeeID, newAge)
	if err != nil {
		log.Fatal(err)
	}

	updatedEmployee, err := getEmployeeByID(db, employeeID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated Employee:")
	fmt.Printf("ID: %d, Name: %s, Age: %d, Salary: %d\n", updatedEmployee.ID, updatedEmployee.Name, updatedEmployee.Age, updatedEmployee.Salary)

	err = deleteEmployee(db, 1)
	if err != nil {
		log.Fatal(err)
	}

	remainingEmployees, err := getEmployees(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Remaining Employees:")
	for _, emp := range remainingEmployees {
		fmt.Printf("ID: %d, Name: %s, Age: %d, Salary: %d\n", emp.ID, emp.Name, emp.Age, emp.Salary)
	}
}

func createEmployee(db *sql.DB, emp employee) error {
	insertQuery := "INSERT INTO employe (name, age, salary) VALUES (?, ?, ?);"
	_, err := db.Exec(insertQuery, emp.Name, emp.Age, emp.Salary)
	return err
}

func getEmployees(db *sql.DB) ([]employee, error) {
	selectQuery := "SELECT id, name, age, salary FROM employe;"
	rows, err := db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []employee
	for rows.Next() {
		var emp employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.Salary)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	return employees, nil
}

func getEmployeeByID(db *sql.DB, id int) (*employee, error) {
	selectQuery := "SELECT id, name, age, salary FROM employe WHERE id = ?;"
	row := db.QueryRow(selectQuery, id)

	var emp employee
	err := row.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.Salary)
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func updateEmployeeAge(db *sql.DB, id, newAge int) error {
	updateQuery := "UPDATE employe SET age = ? WHERE id = ?;"
	_, err := db.Exec(updateQuery, newAge, id)
	return err
}

func deleteEmployee(db *sql.DB, id int) error {
	deleteQuery := "DELETE FROM employe WHERE id = ?;"
	_, err := db.Exec(deleteQuery, id)
	return err
}
