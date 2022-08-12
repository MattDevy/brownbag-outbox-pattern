package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type PubSubAPI interface {
	PublishTopic(ctx context.Context, topic string, message *pubsub.Message) *PublishResult
}
