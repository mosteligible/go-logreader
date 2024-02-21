package notifier

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mosteligible/go-logreader/client/config"
	"github.com/mosteligible/go-logreader/client/customer"
)

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
	response chan<- error,
) {
	log.Printf("Initialization sending request to: %s", url)
	client := &http.Client{}
	var req *http.Request
	var respond *http.Response
	var err error
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
		response <- err
		return
	}

	SetHeaders(req, headers)
	respond, err = client.Do(req)
	if err != nil {
		log.Fatalf("error sending request: %s\n", err.Error())
		response <- err
		return
	}
	if respond.StatusCode > 399 {
		log.Fatalf("API: <%s> respoded with status code: <%d>", url, respond.StatusCode)
		response <- err
		return
	}
	response <- nil
}

func NotifyService(serviceEndpoint string, cust customer.Customer, headers map[string]string, event string) error {
	// send message to concerned services about updates in
	// clientel they should expect and update states accordingly
	var response error
	responseChan := make(chan error)
	go SendRequest(
		serviceEndpoint,
		map[string]string{"api-key": config.Env.BoxUpdApiKey},
		http.MethodPost,
		&map[string]string{"Id": cust.Id, "Name": cust.Name, "Event": event},
		responseChan,
	)
	response = <-responseChan
	if response != nil {
		return response
	}

	return nil
}
