package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "secondhandnotification@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, os.Getenv("SMTP_API_KEY"), os.Getenv("SMTP_SECRET_KEY"))

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("error sending email:", err)
		return err
	}
	return nil
}
