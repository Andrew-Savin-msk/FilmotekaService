package gomaildealer

import (
	"strings"

	"gopkg.in/gomail.v2"
)

type Dealer struct {
	md        *gomail.Dialer
	mail_body string
}

func New(host string, port int, login, password string, mailBody string) *Dealer {
	return &Dealer{
		md:        gomail.NewDialer(host, port, login, password),
		mail_body: mailBody,
	}
}

func (d *Dealer) Send(recipients []string) ([]string, error) {
	unsended := []string{}
	for _, recepient := range recipients {
		mess := gomail.NewMessage()
		mess.SetHeader("Subject", "Welcome to Our Service!")
		mess.SetBody("text/html", strings.ReplaceAll(string(d.mail_body), "[USER_NAME]", d.md.Username))
		mess.SetAddressHeader("To", recepient, recepient)
		err := d.md.DialAndSend(mess)
		if err != nil {
			unsended = append(unsended, recepient)
		}
	}
	return unsended, nil
}
