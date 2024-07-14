package model

import (
	"log"
	"time"
)

func TestFilm() *Film {
	date, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		log.Fatal(err)
	}
	return &Film{
		Name:      "Test Film",
		Desc:      "This is a test film description which is sufficiently long.",
		Date:      date,
		Assesment: 8.5,
	}
}