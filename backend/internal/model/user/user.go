package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Passwd    string `json:"passwd,omitempty"`
	EncPasswd string `json:"-"`
	IsAdmin   bool   `json:"is_admin"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Passwd, validation.By(requiredIf(u.EncPasswd == "")), validation.Length(3, 20)),
	)
}

func (u *User) CompareHashAndPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.EncPasswd), []byte(password))
}

func (u *User) sanitize() {
	u.Passwd = ""
}

func (u *User) Sequre() error {
	err := error(nil)
	u.EncPasswd, err = encryptPasswd(u.Passwd)
	if err != nil {
		return err
	}

	u.sanitize()

	return nil
}

func encryptPasswd(pwd string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(enc), nil
}