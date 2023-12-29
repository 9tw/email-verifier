package config

import (
	"email_verifier/features/user/domain"
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func SendEmail(subject, HTMLbody string, newUser domain.UserCore) error {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Error("config error :", err.Error())
		return nil
	}

	// sender data
	to := []string{newUser.Username}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// Set up authentication information.
	auth := smtp.PlainAuth("", os.Getenv("Email"), os.Getenv("SMTPpwd"), host)
	msg := []byte(
		"From: " + os.Getenv("Entity") + ": <" + os.Getenv("Email") + ">\r\n" +
			"To: " + newUser.Username + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			HTMLbody)
	err = smtp.SendMail(address, auth, os.Getenv("Email"), to, msg)
	if err != nil {
		return err
	}
	fmt.Println("Check for sent email!")
	return nil
}
