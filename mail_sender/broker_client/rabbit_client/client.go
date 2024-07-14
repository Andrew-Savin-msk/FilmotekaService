package rabbitclient

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	// Connection params
	conn *amqp091.Connection

	// Connection chanel params
	ch *amqp091.Channel

	// Connection queue params
	q *amqp091.Queue

	// Messages params
	msgs <-chan amqp091.Delivery

	// Async params
	notifyCloseChan chan *amqp091.Error
	ctx             context.Context
}

func New(URL string, ctx context.Context) (*Client, error) {
	conn, err := amqp091.Dial(URL)
	if err != nil {
		return nil, err
	}

	notifyChan := conn.NotifyClose(make(chan *amqp091.Error))

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

		ctx:             ctx,
		notifyCloseChan: notifyChan,
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
		select {
		case <-c.ctx.Done():
			c.conn.Close()
			c.ch.Close()
			return
		case <-c.notifyCloseChan:
			c.ch.Close()
			return
		case res <- string(ms.Body):
		}
	}
}
