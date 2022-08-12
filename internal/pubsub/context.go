package pubsub

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	pubsubInstructionKey contextKey = "pubsub-instructions"
)

func InstructionsFromContext(ctx context.Context) []string {
	ins, ok := ctx.Value(pubsubInstructionKey).(string)
	if ok {
		return strings.Split(ins, ",")
	}
	return []string{}
}

func contextFromHeader(ctx context.Context, header http.Header) context.Context {
	return context.WithValue(ctx, pubsubInstructionKey, header.Get(string(pubsubInstructionKey)))
}

func ContextFromMetadata(ctx context.Context, metadata map[string]string) context.Context {
	if metadata == nil {
		return ctx
	}
	return context.WithValue(ctx, pubsubInstructionKey, metadata[string(pubsubInstructionKey)])
}

func PubSubInstructionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := contextFromHeader(c.Request.Context(), c.Request.Header)
		req := c.Request.WithContext(ctx)
		c.Request = req
		c.Next()
	}
}
