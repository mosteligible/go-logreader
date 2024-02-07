package config

import "os"

var APP_PORT = os.Getenv("APP_PORT")
var CLIENT_DATA_API_KEY = os.Getenv("CLIENT_DATA_API_KEY")
var CLIENT_DATA_ENDPOINT = os.Getenv("CLIENT_DATA_ENDPOINT")

var BROKER_CONN_URL = os.Getenv("BROKER_CONN_URL")
