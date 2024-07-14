package model_test

import (
	"testing"

	model "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {
	test_cases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				u := model.TestUser()
				return u
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *model.User {
				u := model.TestUser()
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				u := model.TestUser()
				u.Passwd = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *model.User {
				u := model.TestUser()
				u.Passwd = "1"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password, but with encryptet password",
			u: func() *model.User {
				u := model.TestUser()
				u.Passwd = ""
				u.EncPasswd = "encrypted"
				return u
			},
			isValid: true,
		},
	}

	for _, tc := range test_cases {
		if tc.isValid {
			assert.NoError(t, tc.u().Validate(), tc.name)
		} else {
			assert.Error(t, tc.u().Validate(), tc.name)
		}
	}
}
