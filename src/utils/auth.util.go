package utils

import (
	"caapp-server/src/utils/helper"
	"math/rand"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

const letterBytes = "1234567890"

func GenerateRandomCode(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func SendEmail(to, subject, body string) error {
	mailer := gomail.NewMessage()

	mailer.SetHeader("From", os.Getenv("MAIL_USERNAME"))

	mailer.SetHeader("To", to)

	mailer.SetHeader("Subject", subject)

	mailer.SetBody("text/plain", body)

	// Configure the SMTP connection
	dialer := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		int(helper.StringToUInt(os.Getenv("MAIL_PORT"))),
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
	)

	// Send the email
	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
