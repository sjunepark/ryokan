package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"net/smtp"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file: ", "error", err)
	}
	from := os.Getenv("GOOGLE_EMAIL")
	password := os.Getenv("GOOGLE_PASSWORD")
	to := []string{
		"sejun.park@pwc.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("This is a test email message.")
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		slog.Error("Error sending email: ", "error", err)
		return
	}
}
