FROM golang:1.22-alpine3.20

WORKDIR /go/mail_sender

COPY . .

RUN go mod tidy

CMD go run cmd/main.go
