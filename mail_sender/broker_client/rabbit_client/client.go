package rabbitclient

import "github.com/rabbitmq/amqp091-go"

type Client struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
	q    *amqp091.Queue
	msgs <-chan amqp091.Delivery
}

func New(URL string) (*Client, error) {
	conn, err := amqp091.Dial(URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"Temporary",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		ch:   ch,
		q:    &q,
		msgs: msgs,
	}, nil
}

// TODO: Think if we need do msg channel a part of Client structure
func (c *Client) GetMessages() <-chan string {
	msg := make(chan string)
	go c.messagesConvertor(msg)
	return msg
}

func (c *Client) messagesConvertor(res chan string) {
	defer close(res)
	for ms := range c.msgs {
		res <- string(ms.Body)
	}
}
