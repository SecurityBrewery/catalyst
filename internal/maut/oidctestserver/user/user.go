package user

import "golang.org/x/text/language"

type User struct {
	ID                string
	Username          string
	Password          string
	Firstname         string
	Lastname          string
	Email             string
	EmailVerified     bool
	Phone             string
	PhoneVerified     bool
	PreferredLanguage language.Tag
}
