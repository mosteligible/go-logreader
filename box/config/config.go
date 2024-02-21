package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	AppPort       string
	BrokerConnUrl string
	ClientApiKey  string
	ClienUrl      string
}

func newEnvironment() Environment {
	e := Environment{}
	e.Load()
	return e
}

func (e *Environment) Load() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading env vars from .env: %s", err.Error())
	}
	e.AppPort = os.Getenv("APP_PORT")
	e.BrokerConnUrl = os.Getenv("BROKER_CONN_URL")
	e.ClientApiKey = os.Getenv("CLIENT_API_KEY")
	e.ClienUrl = os.Getenv("CLIENT_URL")
	log.Println("broker url:", e.BrokerConnUrl)
	log.Println("app port:", e.AppPort)
}

var Env = newEnvironment()
