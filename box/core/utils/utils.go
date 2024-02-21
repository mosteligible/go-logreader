package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mosteligible/go-logreader/box/config"
	"github.com/mosteligible/go-logreader/box/core/models"
)

func RespondWithJson(w http.ResponseWriter, code int, message interface{}) {
	response, _ := json.Marshal(message)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string, json bool) {
	if !json {
		w.WriteHeader(code)
		w.Write([]byte(message))
	} else {
		RespondWithJson(w, code, map[string]string{"detail": message})
	}
}

func SetHeaders(req *http.Request, headers map[string]string) *http.Request {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req
}

func SendRequest(
	url string,
	headers map[string]string,
	method string,
	postBody *map[string]string,
) (*http.Response, error) {
	log.Printf("Initialization sending request to: %s", url)
	client := &http.Client{}
	var req *http.Request
	var response *http.Response
	var err error
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, url, nil)
	case http.MethodPost:
		var pb []byte
		pb, _ = json.Marshal(postBody)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(pb))
	default:
		panic("Only Supports GET or POST Requests!")
	}
	if err != nil {
		log.Fatalf("error building request: %s\n", err.Error())
		return nil, err
	}

	SetHeaders(req, headers)
	response, err = client.Do(req)
	if err != nil {
		log.Fatalf("error sending request: %s\n", err.Error())
		return nil, err
	}
	if response.StatusCode > 399 {
		log.Fatalf("API: <%s> respoded with status code: <%d>", url, response.StatusCode)
		return nil, err
	}
	return response, err
}

func ParseToStruct[T any](body io.Reader) []T {
	var t []T
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&t); err != nil {
		log.Fatalf("error decoding body into target")
	}
	return t
}

func GetAllClients() []models.Customer {
	var clients []models.Customer
	headers := map[string]string{
		"api-key": config.Env.ClientApiKey,
	}
	resp, err := SendRequest(
		config.Env.ClienUrl, headers, http.MethodGet, nil,
	)
	if err != nil {
		log.Fatalf("error fetching clientel list from %s", config.Env.ClienUrl)
	}

	clients = ParseToStruct[models.Customer](resp.Body)
	fmt.Println("parse from json clients:\n", clients)

	return clients
}
