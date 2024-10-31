package rabbitclient

import (
	"strings"

	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/config"
	"github.com/google/uuid"
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

	// Info fields
	logger *logrus.Entry
}

func New(cfg config.Broker, logger *logrus.Entry) (*Client, error) {
	URL := "amqp://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port[strings.Index(cfg.Port, ":")+1:] + "/"

	conn, err := amqp091.Dial(URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"emails",
		false,
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

		logger: logger,
	}, nil
}

func (c *Client) SendEMailAddreas(addr string) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	err = c.ch.Publish(
		"",
		c.q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(addr),
			MessageId:   uuid.String(),
		},
	)
	if err != nil {
		return err
	}
	c.logger.Infof("send email addreas whith uuid: %s", uuid.String())

	return nil
}
