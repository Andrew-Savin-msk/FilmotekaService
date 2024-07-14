package maildealer

type MailDealer interface {
	Send([]string) ([]string, error)
}
