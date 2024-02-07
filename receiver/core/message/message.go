package message

import (
	"fmt"
	"sync"

	"github.com/mosteligible/go-logreader/receiver/config"
	"github.com/mosteligible/go-logreader/receiver/core/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type LogSender struct {
	ClientId string
	Conn     *amqp.Connection
	mu       sync.Mutex
}

func NewLogSender(clientId string) LogSender {
	conn, err := amqp.Dial(config.BROKER_CONN_URL)
	utils.LogFatalOnError(
		fmt.Sprintf("While connecting to broker for client id: %s", clientId),
		err,
	)
	return LogSender{ClientId: clientId, Conn: conn}
}

func (sender *LogSender) initialize() {
	ch, err := sender.Conn.Channel()
	utils.LogFatalOnError("Failed to open channel.", err)

	err = ch.ExchangeDeclare(
		sender.ClientId,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.LogFatalOnError("Failed to declare exchange!", err)
}
