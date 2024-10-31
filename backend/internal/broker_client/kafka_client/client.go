package kafkaclient

import (
	"encoding/json"
	"strings"

	brockerclient "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/broker_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/config"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Client struct {
	topic    string
	producer sarama.SyncProducer

	logger *logrus.Entry
}

func New(cfg config.Broker, logger *logrus.Entry) (*Client, error) {

	URL := cfg.Host + ":" + cfg.Port[strings.Index(cfg.Port, ":")+1:]

	producer, err := sarama.NewSyncProducer([]string{URL}, nil) // Не река
	if err != nil {
		return nil, err
	}

	return &Client{
		topic:    cfg.Topic,
		producer: producer,
		logger:   logger,
	}, nil
}

func (c *Client) SendEMailAddreas(addr string) error {
	message := brockerclient.Message{
		UUID:  uuid.NewString(),
		Email: addr,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Printf("error marshalling message. error: %s", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: c.topic,
		Value: sarama.ByteEncoder(messageBytes),
	}

	_, _, err = c.producer.SendMessage(msg)
	if err != nil {
		c.logger.Printf("unable to send message. error: %s", err)
		return err
	}

	return nil
}
