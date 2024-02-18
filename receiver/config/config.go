package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Conf struct {
	AppPort            string
	ClientDataApiKey   string
	ClientDataEndpoint string
	BrokerConnUrl      string
}

func newConf() Conf {
	c := Conf{}
	c.Load()
	return c
}

func (c *Conf) Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file.")
	}
	c.AppPort = os.Getenv("APP_PORT")
	c.ClientDataApiKey = os.Getenv("CLIENT_DATA_API_KEY")
	c.ClientDataEndpoint = os.Getenv("CLIENT_DATA_ENDPOINT")
	c.BrokerConnUrl = os.Getenv("BROKER_CONN_URL")
}

var Env = newConf()
