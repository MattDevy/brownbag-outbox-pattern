package events

import "github.com/MattDevy/brownbag-outbox-pattern/internal/models"

const (
	ItemCreatedType string = "com.synack.brownbag.items.created"
)

type ItemCreatedEvent struct {
	Item models.Item
}
