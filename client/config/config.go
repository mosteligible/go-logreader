package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	AppPort             string
	ClientDbUsername    string
	ClientDbPassword    string
	ClientDbHost        string
	ClientDbPort        string
	ClientDbTableName   string
	DbSslMode           string
	SslOk               bool
	PostgresDb          string
	ProjDir             string
	CustomerReadApiKey  string
	ReceiverSvcEndpoint string
	ReceiverUpdApiKey   string
	BoxSvcEndpoint      string
	BoxUpdApiKey        string
}

func newEnvironment() Environment {
	e := Environment{}
	e.Load()
	return e
}

func (e *Environment) Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error trying to read from .env file, continuing with system env vars..")
	}
	e.AppPort = os.Getenv("APP_PORT")
	e.ClientDbUsername = os.Getenv("CLIENT_DB_USERNAME")
	e.ClientDbPassword = os.Getenv("CLIENT_DB_PASSWORD")
	e.PostgresDb = os.Getenv("POSTGRES_DB")
	e.ClientDbHost = os.Getenv("CLIENT_DB_HOST")
	e.ClientDbPort = os.Getenv("CLIENT_DB_PORT")
	e.ClientDbTableName = os.Getenv("CLIENT_DB_TABLE_NAME")
	e.DbSslMode, e.SslOk = os.LookupEnv("DB_SSL_MODE")
	e.ProjDir, _ = os.Getwd()
	e.CustomerReadApiKey = os.Getenv("CUSTOMER_READ_API_KEY")
	e.ReceiverSvcEndpoint = os.Getenv("RECEIVER_SVC_ENDPOINT")
	e.BoxSvcEndpoint = os.Getenv("BOX_SVC_ENDPOINT")
}

var Env = newEnvironment()
