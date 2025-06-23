package entity

import (
	"github.com/google/uuid"
)

// Person is an entity that represents a person. It has a unique identifier.
type Person struct {
	ID   uuid.UUID
	Name string
	Age  int
}
