package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/MattDevy/brownbag-outbox-pattern/internal/events"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/models"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/store"

	"cloud.google.com/go/pubsub"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func (s Server) BeforePostItems(c *gin.Context) {
	// Get item from request
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// set ID (there's a better way to do this with gorm)
	item.ID = uuid.New()

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

	// Publish BEFORE inserting to DB
	res := s.PubSubAPI.PublishTopic(c.Request.Context(), OutputTopic, &pubsub.Message{
		Data:       eb,
		Attributes: s.PubSubAttrs,
	})

	_, err = res.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sItem := &store.Item{
		ID:    item.ID,
		Name:  item.Name,
		Count: item.Count,
		Price: item.Price,
	}
	if err := s.store.CreateItem(c.Request.Context(), sItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sItem)
	return
}
