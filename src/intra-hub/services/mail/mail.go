package mail
import (
    "intra-hub/models"
    "net/smtp"
)

var (
    auth = smtp.PlainAuth("", username, password, host)
)

func SendUserCreated(user *models.User) error {
    return nil
}