package models

import "github.com/google/uuid"

type Item struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name" binding:"required"`
	Count int       `json:"count"`
	Price float32   `json:"price" binding:"required"`
}
