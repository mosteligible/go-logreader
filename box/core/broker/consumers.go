package broker

type Consumers struct{}

func (c *Consumers) Listen() {
	for _, conn := range ConnectionPool {
		go conn.Consume()
	}
}
