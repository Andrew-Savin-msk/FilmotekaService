package rabbitclient

import (
	"context"

	brokerclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/broker_client"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
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

	// Info fields
	logger *logrus.Entry
}

func New(URL string, ctx context.Context, logger *logrus.Entry) (*Client, error) {
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

		logger: logger,
	}, nil
}

// TODO: Think if we need do msg channel a part of Client structure
func (c *Client) GetMessages() <-chan brokerclient.Message {
	msg := make(chan brokerclient.Message)
	go c.messagesConvertor(msg)
	go c.messagesConvertor(msg)
	go c.messagesConvertor(msg)
	return msg
}

func (c *Client) messagesConvertor(res chan brokerclient.Message) {
	defer close(res)
	for ms := range c.msgs {
		logger := c.logger.WithField("message_uuid", ms.MessageId)
		dto := brokerclient.Message{
			UUID: ms.MessageId,
			Mail: string(ms.Body),
		}
		select {
		case <-c.ctx.Done():
			c.conn.Close()
			c.ch.Close()
			return
		case <-c.notifyCloseChan:
			c.logger.Info("lost connection with rabbit client, shoting down the server")
			c.ch.Close()
			return
		case res <- dto:
			logger.Info("read message from queue")
		}
	}
}
