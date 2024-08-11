package utils

import (
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

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

	mailer.SetHeader("From", "duongminhtruong2002.lequydon@gmail.com")

	mailer.SetHeader("To", to)

	mailer.SetHeader("Subject", subject)

	mailer.SetBody("text/plain", body)

	dialer := gomail.NewDialer("smtp.example.com", 587, "duongminhtruong2002.lequydon@gmail.com", "jhda naqz lyrp eozp")

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
