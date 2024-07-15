package maildealer

type MailDealer interface {
	Send(string) error
}
