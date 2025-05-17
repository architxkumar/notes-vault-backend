package utils

import "net/mail"

func EmailValidator(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
