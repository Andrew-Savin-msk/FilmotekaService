package kafkaclient

import (
	"context"
	"strings"

	brokerclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/broker_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/config"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Client struct {
	topic    string
	consumer sarama.Consumer

	logger *logrus.Entry

	ctx context.Context
}

func New(Bc config.Broker, ctx context.Context, logger *logrus.Entry) (*Client, error) {

	URL := Bc.Host + ":" + Bc.Port[strings.Index(Bc.Port, ":")+1:]

	consumer, err := sarama.NewConsumer([]string{URL}, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		topic:    Bc.Topic,
		consumer: consumer,
		logger:   logger,
	}, nil
}

func (c *Client) GetMessages() <-chan brokerclient.Message {
	msg := make(chan brokerclient.Message)
	go c.messagesConvertor(msg)
	go c.messagesConvertor(msg)
	go c.messagesConvertor(msg)
	return msg
}

func (c *Client) messagesConvertor(chan brokerclient.Message) {

}
