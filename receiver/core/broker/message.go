package broker

import (
	"sync"
)

type LogSender struct {
	ClientId string
	Conn     *Connection
	mu       sync.Mutex
}

func NewLogSender(clientId string) LogSender {
	conn := GetConnection(clientId)

	return LogSender{ClientId: clientId, Conn: conn, mu: sync.Mutex{}}
}
