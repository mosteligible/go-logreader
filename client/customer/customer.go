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
		config.Env.ClientDbTableName,
	)
	err := db.QueryRow(
		query,
		cust.Id, cust.Name, cust.Plan, cust.token,
	).Scan(&cust.Id, &cust.Name, &cust.Plan)
	log.Println("config.CLIENT_DB_TABLE_NAME:", config.Env.ClientDbTableName)
	return err
}

func (cust *Customer) UpdateCustomer(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE %s set name=$1, plan=$2, token=$3", config.Env.ClientDbTableName)
	_, err := db.Exec(query, cust.Name, cust.Plan, cust.token)

	return err
}

func (cust *Customer) DeleteCustomer(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", config.Env.ClientDbTableName)
	_, err := db.Query(query, cust.Id)

	return err
}
