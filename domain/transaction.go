package domain

import (
	"time"

	"github.com/google/uuid"
)

// Transaction is a value object that represents a transaction. It has no identifier
type Transaction struct {
	from      uuid.UUID
	amount    float64
	to        uuid.UUID
	createdAt time.Time
}
