package main

import (
	"log"

	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/apiserver"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/config"
)

func main() {
	cfg := config.Load()
	// fmt.Println(cfg)
	err := apiserver.Start(cfg)
	if err != nil {
		log.Fatalf("unable to start server. ended with error: %v", err)
	}
}
