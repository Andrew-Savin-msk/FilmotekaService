package kafkaclient

import (
	"context"
	"encoding/json"
	"strings"

	brokerclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/broker_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/config"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Client struct {
	topic    string
	conn     sarama.Consumer
	consumer sarama.PartitionConsumer

	brokerMessages <-chan *sarama.ConsumerMessage

	logger *logrus.Entry

	ctx context.Context
}

func New(Bc config.Broker, ctx context.Context, logger *logrus.Entry) (*Client, error) {

	URL := Bc.Host + ":" + Bc.Port[strings.Index(Bc.Port, ":")+1:]

	conn, err := sarama.NewConsumer([]string{URL}, nil)
	if err != nil {
		return nil, err
	}

	consumer, err := conn.ConsumePartition(Bc.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}

	return &Client{
		topic:          Bc.Topic,
		conn:           conn,
		consumer:       consumer,
		brokerMessages: consumer.Messages(),
		logger:         logger,
		ctx:            ctx,
	}, nil
}

func (c *Client) GetMessages() <-chan brokerclient.Message {
	msg := make(chan brokerclient.Message)
	go c.messagesConvertor(msg)
	go c.messagesConvertor(msg)
	go c.messagesConvertor(msg)
	return msg
}

func (c *Client) messagesConvertor(output chan brokerclient.Message) {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.brokerMessages:
			res := brokerclient.Message{}

			err := json.Unmarshal(msg.Value, &res)
			if err != nil {
				c.logger.Errorf("unable to unmarshal data error: %s timestamp: %v", err, msg.Timestamp)
			} else {
				c.logger.Errorf("successfully handled message with timestamp: %v", msg.Timestamp)
				output <- res
			}
		}
	}
}
