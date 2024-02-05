package config

import "os"

var APP_PORT = os.Getenv("APP_PORT")
var CLIENT_DB_USERNAME = os.Getenv("CLIENT_DB_USERNAME")
var CLIENT_DB_PASSWORD = os.Getenv("CLIENT_DB_PASSWORD")
var POSTGRES_DB = os.Getenv("POSTGRES_DB")
var CLIENT_DB_HOST = os.Getenv("CLIENT_DB_HOST")
var CLIENT_DB_PORT = os.Getenv("CLIENT_DB_PORT")
var CLIENT_DB_TABLE_NAME = os.Getenv("CLIENT_DB_TABLE_NAME")
var DB_SSL_MODE, SslOk = os.LookupEnv("DB_SSL_MODE")
var PROJ_DIR, _ = os.Getwd()
var CUSTOMER_READ_API_KEY = os.Getenv("CUSTOMER_READ_API_KEY")
