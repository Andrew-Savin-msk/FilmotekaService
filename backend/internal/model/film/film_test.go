package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
)

func TestFilmValidate(t *testing.T) {
	testCases := []struct {
		name    string
		f       func() *film.Film
		isValid bool
	}{
		{
			name: "valid",
			f: func() *film.Film {
				f := film.TestFilm()
				return f
			},
			isValid: true,
		},
		{
			name: "empty name",
			f: func() *film.Film {
				f := film.TestFilm()
				f.Name = ""
				return f
			},
			isValid: false,
		},
		{
			name: "short description",
			f: func() *film.Film {
				f := film.TestFilm()
				f.Desc = "short"
				return f
			},
			isValid: false,
		},
		{
			name: "default date",
			f: func() *film.Film {
				f := film.TestFilm()
				f.Date = time.Time{}
				return f
			},
			isValid: false,
		},
		{
			name: "future date",
			f: func() *film.Film {
				f := film.TestFilm()
				f.Date = time.Now().AddDate(0, 0, 1)
				return f
			},
			isValid: false,
		},
		{
			name: "negative assessment",
			f: func() *film.Film {
				f := film.TestFilm()
				f.Assesment = -1
				return f
			},
			isValid: false,
		},
		{
			name: "too high assessment",
			f: func() *film.Film {
				f := film.TestFilm()
				f.Assesment = 11
				return f
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		if tc.isValid {
			assert.NoError(t, tc.f().Validate(), tc.name)
		} else {
			assert.Error(t, tc.f().Validate(), tc.name)
		}
	}
}
