package messaging

import (
	"context"

	"cloud.google.com/go/pubsub"
)

const (
	_pubSub = "pubsub"
)

// Message struct holds the information which is required to send message to
// messaging services like pub/sub or kafka
type Message struct {
	Project     string
	TopicName   string
	messageType string
	topic       *pubsub.Topic
	ctx         context.Context
}

// NewPubSub function creates an instance of message which will send the data
// to service using the google cloud pub sub library
func NewPubSub(project, topic string) (*Message, error) {
	var m = &Message{Project: project, TopicName: topic}
	var ctx = context.Background()
	var client, err = pubsub.NewClient(ctx, m.Project)
	if err != nil {
		return nil, err
	}
	t, err := client.CreateTopic(ctx, m.TopicName)
	if err != nil {
		return nil, err
	}
	m.topic = t
	m.ctx = ctx
	m.messageType = _pubSub

	return m, nil
}

// Send will check whether message delivery was acknowledged by the service
func (m *Message) Send(msg []byte) bool {
	switch m.messageType {
	case _pubSub:
		var result = m.topic.Publish(m.ctx, &pubsub.Message{
			Data: msg,
		})
		var _, err = result.Get(m.ctx)
		// TODO: may be retry sending the message if it failed?
		return err != nil
	}
	return false
}
