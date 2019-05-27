package mail

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go"
)

type Config struct {
	Key, Domain string
}

type Params struct {
	Sender, Subject string
	Body, Recipient string
	CC, BCC         []string // CC emails
	Timeout         int      // timeout in seconds
}

// SendViaMailgun will try to send the mail using mailgun
func SendViaMailgun(conf *Config, params *Params) (string, string, error) {
	mg := mailgun.NewMailgun(conf.Domain, conf.Key)
	message := mg.NewMessage(params.Sender, params.Subject, params.Body, params.Recipient)
	message.SetHtml(params.Body)

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
