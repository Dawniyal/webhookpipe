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

type UpdateEventPayload struct {
	ID         uuid.UUID       `json:"id" validate:"required,uuid"`
	EndpointID *string         `json:"endpointId"`
	Payload    *map[string]any `json:"payload"`
	Status     *EventStatus    `json:"status" validate:"omitempty,oneof=received pending delivered failed"`
	Active     *bool           `json:"active"`
}

func (p *UpdateEventPayload) Validate() error {
	return validate.Struct(p)
}

type DeleteEventPayload struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (p *DeleteEventPayload) Validate() error {
	return validate.Struct(p)
}
