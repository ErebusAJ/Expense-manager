package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error{
	godotenv.Load()
	
	email := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")
	if email == "" || pass == ""{
		return fmt.Errorf("error verifying email details")
	}

	m := gomail.NewMessage()

	m.SetAddressHeader("From", email, "Budget Breezer")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, email, pass)
	err := d.DialAndSend(m)
	if err != nil{
		return err
	}

	return nil
}