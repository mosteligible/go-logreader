package config

import "os"

var CLIENT_DATA_API_KEY = os.Getenv("CLIENT_DATA_API_KEY")
var CLIENT_DATA_ENDPOINT = os.Getenv("CLIENT_DATA_ENDPOINT")
