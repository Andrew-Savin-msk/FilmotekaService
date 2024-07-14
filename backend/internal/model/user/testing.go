package model

func TestUser() *User {
	return &User{
		Email:  "yandex@mail.ru",
		Passwd: "qwerty",
	}
}
