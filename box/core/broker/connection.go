package broker

import (
	"fmt"
	"log"
	"time"

	"github.com/mosteligible/go-logreader/box/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	ClientId           string
	Conn               *amqp.Connection
	Channel            *amqp.Channel
	Exchange           string
	Queues             []string
	err                chan error
	foreverConsumption chan bool
}

var ConnectionPool = map[string]*Connection{}

func NewConnection(name, exchange string, queues []string) *Connection {
	if conn, ok := ConnectionPool[name]; ok {
		return conn
	}
	if len(queues) == 0 {
		queues = []string{""}
	}
	conn := &Connection{
		ClientId:           name,
		Exchange:           exchange,
		Queues:             queues,
		err:                make(chan error),
		foreverConsumption: make(chan bool),
	}
	conn.Connect()
	ConnectionPool[name] = conn
	return conn
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
	go c.HandleDisconnect()
	go c.listenConnCloseNotif()

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

func (c *Connection) listenConnCloseNotif() {
	<-c.Conn.NotifyClose(make(chan *amqp.Error))
	log.Printf(" [x] Clientid: %s - Connection received closed notification!!", c.ClientId)
	c.foreverConsumption <- true
	c.err <- fmt.Errorf("connection closed: %s", c.ClientId)
	log.Println("Errors propagated to err channels")
}

func (c *Connection) BindQueue() error {
	for _, q := range c.Queues {
		if _, err := c.Channel.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return fmt.Errorf("error in declaring the queue %s", err)
		}
		if err := c.Channel.QueueBind(q, "#", c.ClientId, false, nil); err != nil {
			return fmt.Errorf("queue bind error: %s", err)
		}
	}

	return nil
}

func (c *Connection) Consume() {
	for _, q := range c.Queues {
		err := c.BindQueue()
		if err != nil {
			log.Printf("Error binding queue; %s\n", err.Error())
			return
		}
		deliveries, err := c.Channel.Consume(
			q, "", false, false, false, false, nil,
		)
		if err != nil {
			log.Println("Error in consumption of deliveries..")
			return
		}
		go func() {
			for msg := range deliveries {
				log.Printf(" [x] %s - Received - %s", c.ClientId, msg.Body)
			}
		}()
	}
	<-c.foreverConsumption
	log.Println(" [x] Closing consumer after connection close notification")
}

func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}
	if err := c.BindQueue(); err != nil {
		return err
	}
	return nil
}

func (c *Connection) HandleDisconnect() {
	for {
		time.Sleep(1 * time.Second)

		if err := <-c.err; err != nil {
			log.Printf("Broker disconnect notification received: %s", err.Error())
			for {
				if err := c.Reconnect(); err != nil {
					log.Printf(" [x] error reconnecting to broker: %s", err.Error())
					time.Sleep(1 * time.Second)
					continue
				}
				break
			}
			log.Printf("Restarting consumer for clientid: %s", c.ClientId)
			go c.Consume()
		}
	}
}
