package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS customers (id SERIAL PRIMARY KEY, name TEXT, email TEXT, status TEXT);")
	if err != nil {
		log.Fatal("can't create table customers ", err)
	}

	fmt.Println("create table success.")
}

type Connection interface {
	CreateCustomer(customer *Customer) error
	GetById(id int) (*Customer, error)
	GetAll() ([]*Customer, error)
	UpdateById(id int, customer *Customer) error
	DeleteById(id int) error
}

type connection struct {
	conn *sql.DB
}

func Conn() Connection {
	return &connection{
		conn: db,
	}
}

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (c *connection) CreateCustomer(customer *Customer) error {
	row := c.conn.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", customer.Name, customer.Email, customer.Status)

	err := row.Scan(&customer.ID)
	if err != nil {
		return fmt.Errorf("can't create statement: %w", err)
	}
	return nil
}

func (c *connection) GetById(id int) (*Customer, error) {
	customer := Customer{ID: id}

	stmt, err := c.conn.Prepare("SELECT name, email, status FROM customers where id=$1")
	if err != nil {
		return &customer, fmt.Errorf("can't prepare get statement: %w", err)
	}

	row := stmt.QueryRow(id)

	err = row.Scan(&customer.Name, &customer.Email, &customer.Status)
	if err != nil {
		return &customer, fmt.Errorf("can't scan get statment: %w", err)
	}
	return &customer, nil
}

func (c *connection) GetAll() ([]*Customer, error) {
	var customers []*Customer

	stmt, err := c.conn.Prepare("SELECT * FROM customers")
	if err != nil {
		return customers, fmt.Errorf("can't prepare get statement: %w", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return customers, fmt.Errorf("can't query get statement: %w", err)
	}

	for rows.Next() {
		customer := &Customer{}

		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Status)
		if err != nil {
			return customers, fmt.Errorf("can't scan get statement: %w", err)
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (c *connection) UpdateById(id int, customer *Customer) error {
	stmt, err := c.conn.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		return fmt.Errorf("can't prepare update statement: %w", err)
	}

	if _, err := stmt.Exec(id, &customer.Name, &customer.Email, &customer.Status); err != nil {
		return fmt.Errorf("can't execute update statment: %w", err)
	}
	return nil
}

func (c *connection) DeleteById(id int) error {
	stmt, err := c.conn.Prepare("DELETE FROM customers WHERE name = $1")
	if err != nil {
		return fmt.Errorf("can't prepare delete statement: %w", err)
	}

	if _, err := stmt.Exec(id); err != nil {
		return fmt.Errorf("can't execute delete statment: %w", err)
	}
	return nil
}
