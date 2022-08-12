package server

import (
	"net/http"

	"github.com/MattDevy/brownbag-outbox-pattern/internal/models"
	"github.com/MattDevy/brownbag-outbox-pattern/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s Server) OutboxPostItems(c *gin.Context) {
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
	if err := s.store.CreateItemWithHooks(c.Request.Context(), sItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sItem)
	return
}
