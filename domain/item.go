package entity

import (
	"github.com/google/uuid"
)

// Item is an entity that represents an item. It has a unique identifier.
type Item struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
