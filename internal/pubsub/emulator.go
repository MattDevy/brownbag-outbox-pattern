package pubsub

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

const (
	AttributeKeyEmulator  string = "emulator"
	AttributeValueFail    string = "fail"
	AttributeValueRetry5  string = "retry-5"
	AttributeValueTimeout string = "timeout"
	AttributeValueSucceed string = "succeed"
)

var (
	FailedError error = errors.New("some pubsub failure")
)

func CombineAttributes(attributes ...string) string {
	return strings.Join(attributes, ",")
}

type Emulator struct{}

func NewEmulator() *Emulator {
	return &Emulator{}
}

func (e Emulator) PublishTopic(ctx context.Context, topic string, message *pubsub.Message) *PublishResult {
	var (
		instructions = InstructionsFromContext(ctx)
		res          = NewPublishResult()
	)
	for _, instruction := range instructions {
		switch instruction {
		case AttributeValueFail:
			SetPublishResult(res, "test", FailedError)
			zap.L().Info("message failed", zap.Any("message", message))
			return res
		case AttributeValueRetry5:
			for i := 0; i < 5; i++ {
				// fake io
				<-time.After(1 * time.Second)
				zap.L().Info("retrying", zap.Int("attempts", i+1))
			}
		case AttributeValueTimeout:
			<-ctx.Done()
			SetPublishResult(res, "test", ctx.Err())
			zap.L().Info("message timed out", zap.Any("message", message))
			return res
		case AttributeValueSucceed:
			SetPublishResult(res, "test", nil)
			zap.L().Info("message sent", zap.Any("message", message))
			return res
		default:
			SetPublishResult(res, "test", fmt.Errorf("unknown instruction: %s", instruction))
			return res
		}
	}

	// default is success
	SetPublishResult(res, "test", nil)
	zap.L().Info("message sent", zap.Any("message", message))
	return res
}
