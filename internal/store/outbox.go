package store

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/MattDevy/brownbag-outbox-pattern/internal/events"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/models"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/pubsub"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

const (
	OutboxTopic string = "events"
)

func NewItemCreatedEvent(ctx context.Context, item Item) (*Event, error) {
	instructions := pubsub.InstructionsFromContext(ctx)
	event := cloudevents.NewEvent()
	event.SetID(uuid.NewString())
	event.SetSource("github.com/MattDevy/brownbag-outbox-pattern/internal/store")
	event.SetType(events.ItemCreatedType)
	event.SetData(
		cloudevents.ApplicationJSON,
		events.ItemCreatedEvent{Item: models.Item{
			ID:    item.ID,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		}},
	)
	eb, err := event.MarshalJSON()
	if err != nil {
		return nil, err
	}

	md, err := json.Marshal(map[string]string{
		"pubsub-instructions": strings.Join(instructions, ","),
	})
	if err != nil {
		return nil, err
	}

	return &Event{
		Uuid:     uuid.NewString(),
		Payload:  JSON(eb),
		Metadata: JSON(md),
	}, nil
}

type Event struct {
	Offset    int `gorm:"AUTO_INCREMENT;->"`
	Uuid      string
	CreatedAt time.Time
	Payload   JSON
	Metadata  JSON
}

func (Event) TableName() string {
	return "outbox_events"
}
