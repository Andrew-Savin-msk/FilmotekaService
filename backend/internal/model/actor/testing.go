package model

import (
	"log"
	"time"
)

func TestActor() *Actor {
	birth, err := time.Parse("01-02-2006", "01-01-2001")
	if err != nil {
		log.Fatal(err)
	}
	return &Actor{
		Name:      "John Poul",
		Gen:       "man",
		Birthdate: birth,
	}
}