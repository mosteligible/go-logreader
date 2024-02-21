package customer

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mosteligible/go-logreader/client/config"
)

func GetCustomer(db *sql.DB, id string) (Customer, error) {
	customer := Customer{Id: id}
	query := fmt.Sprintf(
		"SELECT id, name, plan from %s where id = $1",
		config.Env.ClientDbTableName,
	)
	row := db.QueryRow(
		query, id,
	)

	err := row.Scan(&customer.Id, &customer.Name, &customer.Plan)
	if err != nil {
		log.Printf("Error getting customer for id: %s - error: %s", id, err.Error())
		return customer, err
	}

	return customer, nil
}

func GetAllCustomers(db *sql.DB) ([]Customer, error) {
	var allCustomers = []Customer{}

	query := fmt.Sprintf("SELECT * from %s", config.Env.ClientDbTableName)
	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c Customer
		if err := rows.Scan(
			&c.Id, &c.Name, &c.token, &c.Plan,
		); err != nil {
			return nil, err
		}
		allCustomers = append(allCustomers, c)
	}

	return allCustomers, nil
}
