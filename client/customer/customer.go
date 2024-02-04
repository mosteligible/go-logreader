package customer

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	"github.com/mosteligible/go-logreader/client/config"
)

type Customer struct {
	Name  string
	Id    string
	Ip    string
	Plan  string
	token string
}

func (cust *Customer) GetToken() string {
	return cust.token
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (cust *Customer) AddCustomer(db *sql.DB) error {
	cust.token = generateToken()
	query := fmt.Sprintf(
		"INSERT INTO %s(id, name, plan, token) VALUES($1, $2, $3, $4) RETURNING id, name, plan",
		config.CLIENT_DB_TABLE_NAME,
	)
	err := db.QueryRow(
		query,
		cust.Id, cust.Name, cust.Plan, cust.token,
	).Scan(&cust.Id, &cust.Name, &cust.Plan)
	log.Println("config.CLIENT_DB_TABLE_NAME:", config.CLIENT_DB_TABLE_NAME)
	return err
}

func (cust *Customer) UpdateCustomer(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE ? set name=?, plan=?, token=?",
		config.CLIENT_DB_TABLE_NAME, cust.Name, cust.Plan, cust.token,
	)

	return err
}

func (cust *Customer) RemoveCustomer(db *sql.DB) error {
	_, err := db.Query("DELETE FROM %s WHERE id=? LIMIT 1", cust.Id)

	return err
}
