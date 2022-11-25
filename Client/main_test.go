package main_test

import (
	"log"
	"os"
	"testing"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS clients
(
    id SERIAL,
    name TEXT NOT NULL,
    plan TEXT NOT NULL,
    CONSTRAINT clients_pkey PRIMARY KEY (id)
)`

var a main.App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)

}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}
