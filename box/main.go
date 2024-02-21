package main

import (
	"github.com/mosteligible/go-logreader/box/app"
	"github.com/mosteligible/go-logreader/box/core/broker"
	"github.com/mosteligible/go-logreader/box/core/utils"
)

func startConsumers(sigConsumerStart chan<- bool) {
	clients := utils.GetAllClients()
	for _, cust := range clients {
		conn := broker.NewConnection(cust.Id, cust.Id, []string{})
		go conn.Consume()
	}
	sigConsumerStart <- true
}

func main() {
	sigConsumerStart := make(chan bool)
	go startConsumers(sigConsumerStart)
	<-sigConsumerStart
	close(sigConsumerStart)
	app := app.NewApp()
	app.Run()
}
