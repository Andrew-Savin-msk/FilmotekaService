package brokerclient

type Client interface {
	GetMessages() <-chan string
}
