package event

import (
	"time"

	"github.com/google/uuid"
)

type Event_Status string

const (
	StatusReceived  Event_Status = "received"
	StatusPending   Event_Status = "pending"
	StatusDelivered Event_Status = "delivered"
	StatusFailed    Event_Status = "failed"
)

type Event struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	EndpointID string         `json:"endpointId" db:"endpoint_id"`
	Payload    map[string]any `json:"payload" db:"payload"`
	Status     Event_Status   `json:"status" db:"status"`
	CreatedAt  time.Time      `json:"createdAt" db:"created_at"`
}
