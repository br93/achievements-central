package data

import "net/mail"

type Validation struct{}

func (*Validation) ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
