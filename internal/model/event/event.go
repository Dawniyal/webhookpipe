package event

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	StatusReceived  EventStatus = "received"
	StatusPending   EventStatus = "pending"
	StatusDelivered EventStatus = "delivered"
	StatusFailed    EventStatus = "failed"
)

type Event struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	EndpointID string         `json:"endpointId" db:"endpoint_id"`
	Payload    map[string]any `json:"payload" db:"payload"`
	Status     EventStatus    `json:"status" db:"status"`
	CreatedAt  time.Time      `json:"createdAt" db:"created_at"`
}
