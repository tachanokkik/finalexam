package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tachanokkik/finalexam/customer"
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
}

type Connection interface {
	CreateCustomer(customer customer.Customer) error
	GetById(id int) (*customer.Customer, error)
	GetAll() ([]*customer.Customer, error)
	UpdateById(id int, customer customer.Customer) error
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

func (c *connection) CreateCustomer(customer customer.Customer) error {
	row := c.conn.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", customer.Name, customer.Email, customer.Status)

	err := row.Scan(&customer.ID)
	if err != nil {
		return fmt.Errorf("can't create statement: %w", err)
	}
	return nil
}

func (c *connection) GetById(id int) (*customer.Customer, error) {
	customer := customer.Customer{ID: id}

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

func (c *connection) GetAll() ([]*customer.Customer, error) {
	var customers []*customer.Customer

	stmt, err := c.conn.Prepare("SELECT * FROM customers")
	if err != nil {
		return customers, fmt.Errorf("can't prepare get statement: %w", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return customers, fmt.Errorf("can't query get statement: %w", err)
	}

	for rows.Next() {
		customer := &customer.Customer{}

		err := rows.Scan(&customer)
		if err != nil {
			return customers, fmt.Errorf("can't scan get statement: %w", err)
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (c *connection) UpdateById(id int, customer customer.Customer) error {
	stmt, err := c.conn.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		return fmt.Errorf("can't prepare update statement: %w", err)
	}

	if _, err := stmt.Exec(id, customer.Name, customer.Email, customer.Status); err != nil {
		return fmt.Errorf("can't execute update statment: %w", err)
	}
	return nil
}

func (c *connection) DeleteById(id int) error {
	stmt, err := c.conn.Prepare("DELETE FROM customerss WHERE name = $1")
	if err != nil {
		return fmt.Errorf("can't prepare delete statement: %w", err)
	}

	if _, err := stmt.Exec(id); err != nil {
		return fmt.Errorf("can't execute delete statment: %w", err)
	}
	return nil
}
