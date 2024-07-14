package mailserver

import (
	"context"
	"sync"

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

	wg.Add(1)
	go s.messagesSender(&wg, s.bc.GetMessages())
	wg.Wait()
	return nil
}

func (s *server) messagesSender(wg *sync.WaitGroup, input <-chan string) {
	defer wg.Done()
	for user := range input {
		select {
		case <-s.ctx.Done():
			return
		default:
			err := s.md.Send(user)
			if err != nil {

			}
		}
	}
}
