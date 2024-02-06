package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mosteligible/go-logreader/receiver/config"
	"github.com/mosteligible/go-logreader/receiver/core/utils"
)

type Client struct {
	Id   string
	Plan string
	Name string
	Ip   string
}

func (c *Client) Validate() (bool, error) {
	var responseComm = make(chan utils.CommResponse)
	var communicatedResponse utils.CommResponse
	var res *http.Response
	endpoint := fmt.Sprintf("%s/customers/%s", config.CLIENT_DATA_ENDPOINT, c.Id)
	headers := map[string]string{"api-key": config.CLIENT_DATA_API_KEY}

	fmt.Println("Starting goroutine..")
	go utils.SendRequest(endpoint, headers, http.MethodGet, nil, responseComm)
	communicatedResponse = <-responseComm
	if communicatedResponse.Err != nil {
		log.Fatalf("Error reading body: %s\n", communicatedResponse.Err.Error())
		return false, communicatedResponse.Err
	}
	res = communicatedResponse.Response
	defer res.Body.Close()
	defer close(responseComm)

	var obtainedClient Client
	decoder := json.NewDecoder(res.Body)
	fmt.Printf("Response body: %s\n", res.Body)
	if err := decoder.Decode(&obtainedClient); err != nil {
		log.Fatalf("Error reading body: %s\nBody: %s\n", err.Error(), res.Body)
		return false, err
	}
	c.Id = obtainedClient.Id
	c.Name = obtainedClient.Name
	c.Plan = obtainedClient.Plan
	return true, nil
}

func (c *Client) String() string {
	return fmt.Sprintf(
		"Client(id: %s, name: %s, plan: %s)\n",
		c.Id, c.Name, c.Plan,
	)
}
