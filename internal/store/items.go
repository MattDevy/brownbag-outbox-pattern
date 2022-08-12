package store

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID    uuid.UUID `gorm:"primarykey"`
	Name  string
	Count int
	Price float32
}

func (i *Item) AfterCreate(tx *gorm.DB) (err error) {
	e, err := NewItemCreatedEvent(tx.Statement.Context, *i)
	if err != nil {
		return err
	}
	return tx.Model(&Event{}).Create(e).Error
}
