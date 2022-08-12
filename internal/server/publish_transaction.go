package server

import (
	"context"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/events"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/models"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/store"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s Server) TransactionPostItems(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// set ID (there's a better way to do this with gorm)
	item.ID = uuid.New()

	sItem := &store.Item{
		ID:    item.ID,
		Name:  item.Name,
		Count: item.Count,
		Price: item.Price,
	}

	// make event
	event := cloudevents.NewEvent()
	event.SetID(uuid.NewString())
	event.SetSource("github.com/MattDevy/brownbag-outbox-pattern/internal/server")
	event.SetType(events.ItemCreatedType)
	event.SetData(cloudevents.ApplicationJSON, events.ItemCreatedEvent{Item: item})
	eb, err := event.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.store.CreateItemTx(c.Request.Context(), sItem, func(ctx context.Context) error {
		// Publish IN db transaction
		res := s.PubSubAPI.PublishTopic(c.Request.Context(), OutputTopic, &pubsub.Message{
			Data:       eb,
			Attributes: s.PubSubAttrs,
		})

		_, err = res.Get(c.Request.Context())
		return err
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sItem)
	return
}
