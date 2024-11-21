package model

import (
	"time"

	model "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Film struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"description,omitempty"`
	Date      time.Time `json:"release_date"`
	Assesment float32   `json:"assesment"`
}

func (f *Film) Validate() error {
	return validation.ValidateStruct(
		f,
		validation.Field(&f.Name, validation.Required, validation.Length(1, 150)),
		validation.Field(&f.Desc, validation.Required, validation.Length(10, 1000)),
		validation.Field(&f.Date, validation.Required, validation.By(model.IsDateValid())),
		validation.Field(&f.Assesment, validation.Required, validation.Max(10.0), validation.Min(0.0)),
	)
}
