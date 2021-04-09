package mail

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type Config struct {
	Key, Domain string
}

type MailjetConfig struct {
	PubKey, PrivateKey string
}

type Params struct {
	Sender, Subject string
	Body, Recipient string
	ReplyTo         string
	CC, BCC         []string // CC emails
	Timeout         int      // timeout in seconds
}

type MailjetParams struct {
	SenderEmail, SenderName string
	RecipientEmail, RecipientName string
	Subject string
	TextPart, HtmlPart string
}

// SendViaMailgun will try to send the mail using mailgun
func SendViaMailgun(conf *Config, params *Params) (string, string, error) {
	mg := mailgun.NewMailgun(conf.Domain, conf.Key)
	message := mg.NewMessage(params.Sender, params.Subject, params.Body, params.Recipient)
	message.SetHtml(params.Body)
	if len(params.ReplyTo) > 0 {
		message.SetReplyTo(params.ReplyTo)
	}

	for _, emailID := range params.CC {
		message.AddCC(emailID)
	}
	for _, emailID := range params.BCC {
		message.AddBCC(emailID)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	return resp, id, err
}

// SendViaMailjet will try to send the mail using mailjet
func SendViaMailjet(conf *MailjetConfig, params *Params) (string, error) {
	mailjetClient := mailjet.NewMailjetClient(conf.PubKey, conf.PrivateKey)
	messagesInfo := []mailjet.InfoMessagesV31 {
		mailjet.InfoMessagesV31{
		  From: &mailjet.RecipientV31{
			Email: params.SenderEmail,
			Name: params.SenderName,
		  },
		  To: &mailjet.RecipientsV31{
			mailjet.RecipientV31 {
			  Email: params.RecipientEmail,
			  Name: params.RecipientName,
			},
		  },
		  Subject: params.Subject,
		  TextPart: params.TextPart,
		  HTMLPart: params.HTMLPart
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}
