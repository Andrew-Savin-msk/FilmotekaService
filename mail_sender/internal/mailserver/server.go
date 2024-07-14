package mailserver

import (
	"context"

	brokerclient "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/broker_client"
	maildealer "github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/mail_dealer"
	"github.com/sirupsen/logrus"
)

type server struct {
	md     maildealer.MailDealer
	bc     brokerclient.Client
	logger *logrus.Logger

	ctx context.Context
}
