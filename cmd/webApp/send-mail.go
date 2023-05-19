package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
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
	var msgBody string

	if msg.Template == "" {
		msgBody = msg.Content
	} else {
		// Get the file
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", msg.Template))
		if err != nil {
			log.Println("failed to read from file")
		}

		// Replace the placeholder with the data to inject
		mailTemplate := string(data)
		mailToSend := strings.Replace(mailTemplate, "[%body%]", msg.Content, 1)

		// assign the string to msgBody
		msgBody = mailToSend
	}
	email.SetBody(mail.TextHTML, msgBody)

	err = email.Send(client)
	if err != nil {
		log.Println("Email not sent!")
	} else {
		log.Println("Email sent")
	}
}
