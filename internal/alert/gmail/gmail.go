package gmail

import (
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"net/smtp"
	"os"
	"strings"
)

func SendEmails(availableDates []civil.Date) (sent bool) {
	if len(availableDates) == 0 {
		return false
	}

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file: ", "error", err)
		return false
	}

	from := os.Getenv("GOOGLE_EMAIL")
	password := os.Getenv("GOOGLE_PASSWORD")
	to := []string{
		"junepark202012@gmail.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)

	datesString := joinDates(availableDates)

	msgSubject := "Ginzanso available dates"
	msgBody := fmt.Sprintf("Available dates: %s", datesString)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", msgSubject, msgBody))

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		slog.Error("Error sending email: ", "error", err)
		return false
	}

	return true
}

func joinDates(dates []civil.Date) string {
	stringDates := make([]string, len(dates))
	for i, date := range dates {
		stringDates[i] = date.String()
	}
	return strings.Join(stringDates, "; ")
}
