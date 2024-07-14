package model_test

import (
	"testing"
	"time"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	"github.com/stretchr/testify/assert"
)

func TestActorValidate(t *testing.T) {
	test_cases := []struct {
		name    string
		a       func() *actor.Actor
		isValid bool
	}{
		{
			name: "valid",
			a: func() *actor.Actor {
				a := actor.TestActor()
				return a
			},
			isValid: true,
		},
		{
			name: "empty name",
			a: func() *actor.Actor {
				a := actor.TestActor()
				a.Name = ""
				return a
			},
			isValid: false,
		},
		{
			name: "empty password",
			a: func() *actor.Actor {
				a := actor.TestActor()
				a.Gen = ""
				return a
			},
			isValid: false,
		},
		{
			name: "default birthdate",
			a: func() *actor.Actor {
				a := actor.TestActor()
				a.Birthdate = time.Time{}
				return a
			},
			isValid: false,
		},
	}

	for _, tc := range test_cases {
		if tc.isValid {
			assert.NoError(t, tc.a().Validate(), tc.name)
		} else {
			assert.Error(t, tc.a().Validate(), tc.name)
		}
	}
}
