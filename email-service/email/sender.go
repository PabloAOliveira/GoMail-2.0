package email

import (
	"email-service/models"
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(task models.EmailTask) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("configurações SMTP incompletas")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := []byte("To: " + task.To + "\r\n" +
		"Subject: " + task.Subject + "\r\n" +
		"\r\n" +
		task.Body + "\r\n")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, smtpUser, []string{task.To}, msg)
}
