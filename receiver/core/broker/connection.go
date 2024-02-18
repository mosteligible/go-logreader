package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/mosteligible/go-logreader/receiver/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	ClientId string
	Conn     *amqp.Connection
	Channel  *amqp.Channel
	Exchange string
	Queues   []string
	err      chan error
}

func NewConnection(name, exchange string, queues []string) *Connection {
	if c, ok := ConnectionPool[name]; ok {
		return c
	}
	c := &Connection{
		ClientId: name,
		Exchange: exchange,
		Queues:   queues,
	}
	c.Connect()
	ConnectionPool[name] = c
	return c
}

func GetConnection(name string) *Connection {
	if val, ok := ConnectionPool[name]; ok {
		return val
	}
	return NewConnection(name, name, []string{})
}

func (c *Connection) Connect() error {
	var err error
	c.Conn, err = amqp.Dial(config.Env.BrokerConnUrl)
	if err != nil {
		return err
	}
	go func() {
		<-c.Conn.NotifyClose(make(chan *amqp.Error))
		c.err <- fmt.Errorf("Connection closed: %s", c.ClientId)
	}()

	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return err
	}

	if err := c.Channel.ExchangeDeclare(
		c.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	return nil
}

func (c *Connection) BindQueue() error {
	for _, q := range c.Queues {
		if _, err := c.Channel.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return fmt.Errorf("error in declaring the queue %s", err)
		}
		if err := c.Channel.QueueBind(q, c.ClientId, c.ClientId, false, nil); err != nil {
			return fmt.Errorf("queue bind error: %s", err)
		}
	}

	return nil
}

func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}
	return nil
}

func (c *Connection) Send(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case err := <-c.err:
		if err != nil {
			c.Reconnect()
			return nil
		}
	default:
		err := c.Channel.PublishWithContext(
			ctx,        // context
			c.Exchange, // exchange
			c.ClientId, // key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		if err != nil {
			return fmt.Errorf("error in publishing message: %s", err.Error())
		}
	}

	return nil
}

var ConnectionPool = map[string]*Connection{}
