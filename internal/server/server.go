package server

import (
	"github.com/MattDevy/brownbag-outbox-pattern/internal/pubsub"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/store"
	"github.com/gin-gonic/gin"
)

const OutputTopic = "brownbag.items.events"

const (
	BeforeBasePath      string = "/before"
	AfterBasePath       string = "/after"
	TransactionBasePath string = "/transaction"
	OutboxBasePath      string = "/outbox"
)

const (
	ItemsPath string = "/items"
)

type Server struct {
	*gin.Engine
	pubsub.PubSubAPI
	PubSubAttrs map[string]string
	store       store.StoreInterface
}

func NewServer(store store.StoreInterface, ps pubsub.PubSubAPI) *Server {
	server := &Server{
		Engine:      gin.Default(),
		store:       store,
		PubSubAPI:   ps,
		PubSubAttrs: make(map[string]string),
	}

	router := server.Engine
	router.Use(pubsub.PubSubInstructionMiddleware())

	// set of endpoints that publish to pubsub before database operations
	before := router.Group(BeforeBasePath)
	before.POST(ItemsPath, server.BeforePostItems)

	// set of endpoints that publish to pubsub after database operations
	after := router.Group(AfterBasePath)
	after.POST(ItemsPath, server.AfterPostItems)

	// set of endpoints that publish to pubsub in-transaction database operations
	transaction := router.Group(TransactionBasePath)
	transaction.POST(ItemsPath, server.TransactionPostItems)

	// set of endpoints that publish to pubsub with outbox database operations
	outbox := router.Group(OutboxBasePath)
	outbox.POST(ItemsPath, server.OutboxPostItems)

	return server
}
