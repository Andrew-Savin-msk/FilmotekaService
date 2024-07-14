package mailserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	brokerclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/broker_client"
	rabbitclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/broker_client/rabbit_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/config"
	maildealer "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/mail_dealer"
	gomaildealer "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/mail_dealer/go_mail_dealer"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnknownMD = errors.New("unknown type of mail dealer")
	ErrUnknownBC = errors.New("unknown type of broker client")
)

func Start(cfg *config.Config) error {
	log := setLog(cfg.LogLevel)

	body, err := loadBody(cfg.MailBodyPath)
	if err != nil {
		return err
	}

	dealer, err := setMailDealer(cfg.MDType, cfg.Host, cfg.Login, cfg.Password, body)
	if err != nil {
		return err
	}

	ctx := context.Background()

	client, err := setBrokerClient(cfg.BrokerType, cfg.BrokerURL, ctx)
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

	htmlBody, err := os.ReadFile("templates/mail_body.html")
	if err != nil {
		log.Fatalf("unable to load body template, ended with error: %s", err)
	}
	return string(htmlBody), nil
}

func setMailDealer(dealerName string, host, login, password string, mailBody string) (maildealer.MailDealer, error) {
	switch dealerName {
	case "gomail", "go-mail", "go_mail":
		return gomaildealer.New(host, 587, login, password, mailBody), nil
	}
	return nil, ErrUnknownMD
}

func setBrokerClient(name, URL string, ctx context.Context) (brokerclient.Client, error) {
	switch strings.ToLower(name) {
	case "rabbitmq", "rabbit_mq", "rabbit":
		return rabbitclient.New(URL, ctx)
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
