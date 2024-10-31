package mailserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	brokerclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/broker_client"
	rabbitclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/broker_client/rabbit_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/config"
	maildealer "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/mail_dealer"
	gomaildealer "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/mail_dealer/go_mail_dealer"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnknownMD = errors.New("unknown type of mail dealer")
	ErrUnknownBC = errors.New("unknown type of broker client")
)

func Start(cfg *config.Config) error {
	log := setLog(cfg.Srv.LogLevel)

	dealer, err := setMailDealer(cfg.Send)
	if err != nil {
		return err
	}

	ctx := context.Background()

	client, err := setBrokerClient(cfg.Bc, ctx, log)
	if err != nil {
		return err
	}

	srv := newServer(log, dealer, client, ctx)

	err = srv.run()
	if err != nil {
		return err
	}

	return nil

}

func setLog(level string) *logrus.Logger {
	log := logrus.New()
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	}
	fmt.Printf("logger set in level: %s \n", level)
	return log
}

func loadBody(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	htmlBody, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to load body template, ended with error: %s", err)
	}
	return string(htmlBody), nil
}

func setMailDealer(Send config.Sender) (maildealer.MailDealer, error) {
	body, err := loadBody(Send.MailBodyPath)
	if err != nil {
		return nil, err
	}
	switch Send.MDType {
	case "gomail", "go-mail", "go_mail":
		return gomaildealer.New(Send.Host, 587, Send.Login, Send.Password, body), nil
	}
	return nil, ErrUnknownMD
}

func setBrokerClient(Bc config.Broker, ctx context.Context, logger *logrus.Logger) (brokerclient.Client, error) {
	switch strings.ToLower(Bc.BrokerType) {
	case "rabbitmq", "rabbit_mq", "rabbit":
		return rabbitclient.New(Bc, ctx, logrus.NewEntry(logger))
	case "kafka", "apache-kafka", "mannaya":
		return kafkaclient.New(Bc, ctx, logrus.NewEntry(logger))
	}
	return nil, ErrUnknownBC
}

func newServer(log *logrus.Logger, MD maildealer.MailDealer, client brokerclient.Client, ctx context.Context) *server {
	srv := server{
		md:     MD,
		bc:     client,
		logger: log,

		ctx: ctx,
	}
	return &srv
}
