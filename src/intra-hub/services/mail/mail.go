package mail

import (
	"net/smtp"

	"intra-hub/confperso"
	"intra-hub/models"
)

var (
	auth = smtp.PlainAuth("", confperso.EmailUsername, confperso.EmailPassword, confperso.EmailHost)
)

func SendUserCreated(user *models.User) error {
	return nil
}
