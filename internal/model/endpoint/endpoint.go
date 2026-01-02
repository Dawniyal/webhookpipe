package endpoint

import "time"

type Endpoint struct {
	ID        string    `json:"id" db:"id"`
	TargetURL string    `json:"targetUrl" db:"target_url"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
