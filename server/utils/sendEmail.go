package utils

import (
	"fmt"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	From string
	To string
	Subject string
	Text string
	Html string
}

func SendEmail(data EmailData) error {
	// get user and pass
	user := config.AppConfig.Email.User
	pass := config.AppConfig.Email.Pass

	// get sender ready 
	s := gomail.NewMessage()

	s.SetAddressHeader("From", user, "Team Ivents Plannerz🎉")
	s.SetHeader("To", data.To)
	s.SetHeader("Subject", data.Subject)
	s.SetBody("text/plain", data.Text)
	s.AddAlternative("text/html", data.Html)

	// get transporter ready 
	t := gomail.NewDialer("smtp.gmail.com", 465, user, pass)

	// try sending the mail 
	err := t.DialAndSend(s)
	if err != nil {
		fmt.Println("Couldn't send Email")
	}

	fmt.Println("SMS Sent✅")
	return nil
}