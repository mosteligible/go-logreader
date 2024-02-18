package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mosteligible/go-logreader/receiver/core/broker"
)

type CommResponse struct {
	Response *http.Response
	Err      error
}

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

func LogFatalOnError(msg string, err error) {
	if err != nil {
		log.Panicf("Error: %s - %s", err.Error(), msg)
	}
}

func SendRequest(
	url string,
	headers map[string]string,
	method string,
	postBody *map[string]string,
	response chan<- CommResponse,
) {
	log.Printf("Initialization sending request to: %s", url)
	client := &http.Client{}
	var req *http.Request
	var respond *http.Response
	var err error
	res := CommResponse{Err: errors.New("Unable to communicate to api")}
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, url, nil)
	case http.MethodPost:
		var pb []byte
		pb, err = json.Marshal(postBody)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(pb))
	default:
		panic("Only Supports GET or POST Requests!")
	}
	if err != nil {
		log.Fatalf("error building request: %s\n", err.Error())
		res.Err = err
		response <- res
		return
	}

	SetHeaders(req, headers)
	respond, err = client.Do(req)
	if err != nil {
		log.Fatalf("error sending request: %s\n", err.Error())
		res.Err = err
		response <- res
		return
	}
	if respond.StatusCode > 399 {
		log.Fatalf("API: <%s> respoded with status code: <%d>", url, respond.StatusCode)
		response <- res
		return
	}
	res.Response = respond
	res.Err = nil
	response <- res
}

func SendMsgWithRetries(msg string, conn *broker.Connection) error {
	var err error = nil
	for i := 0; i < 10; i++ {
		if err = conn.Send(msg); err != nil {
			continue
		}
		break
	}
	return err
}
