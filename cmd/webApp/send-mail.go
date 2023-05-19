package main

import (
	"log"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func ListenForMail() {
	// Runs infinitely
	go func() {
		for {
			msg := <-app.MailChan // This is how I will talk to this go routine from any other place in the project
			log.Println("Got message from the channel, sending mail")
			sendMail(msg)
		}
	}()

}

func sendMail(msg models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.SendTimeout = time.Second * 10
	server.ConnectTimeout = time.Second * 10
	server.KeepAlive = false
	server.Port = 1025

	client, err := server.Connect()
	if err != nil {
		log.Println("could not connect to server")
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextHTML, msg.Content)

	err = email.Send(client)
	if err != nil {
		log.Println("Email not sent!")
	} else {
		log.Println("Email sent")
	}
}
