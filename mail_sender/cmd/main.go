package main

import (
	"fmt"
	"log"

	"github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/config"
	"github.com/Andrew-Savin-msk/filmoteka-service/mail-sender/internal/mailserver"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg)
	err := mailserver.Start(cfg)
	if err != nil {
		log.Fatalf("unable to start server, ended with error: %s", err)
	}
}
