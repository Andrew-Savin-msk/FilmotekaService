package mailserver

import (
	"context"
	"sync"
	"time"

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

// TODO:
func (s *server) run() error {
	var wg sync.WaitGroup
	s.logger.Infof("Mail service started work")
	wg.Add(1)
	go s.messagesSender(&wg, s.bc.GetMessages())
	wg.Wait()
	return nil
}

func (s *server) messagesSender(wg *sync.WaitGroup, input <-chan brokerclient.Message) {
	defer wg.Done()
	for user := range input {
		logger := s.logger.WithField("message_uuid", user.UUID)
		start := time.Now()
		logger.Infof("received message at sending")
		select {
		case <-s.ctx.Done():
			s.logger.Infof("Mail service stopped work with context value: %s", s.ctx.Err())
			return
		default:
			err := s.md.Send(user.Mail)
			if err != nil {
				logger.Infof("unable to send email due to error: %s", err)
			} else {
				logger.Infof("succesfuly sent in %v", time.Since(start))
			}
		}
	}
}
