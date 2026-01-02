package event

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type CreateEventPayload struct {
	EndpointID string         `json:"endpointId" validate:"required"`
	Payload    map[string]any `json:"payload" validate:"required"`
}

func (p *CreateEventPayload) Validate() error {
	if err := validate.Struct(p); err != nil {
		return err
	}

	if _, err := json.Marshal(p.Payload); err != nil {
		return err
	}

	return nil
}

type GetEventByIDPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *GetEventByIDPayload) Validate() error {
	return validate.Struct(p)
}
