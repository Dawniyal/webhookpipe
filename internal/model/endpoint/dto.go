package endpoint

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type AddEndpointPayload struct {
	ID        string `json:"id" validate:"required"`
	TargetURL string `json:"targetUrl" validate:"required,url"`
	Active    bool   `json:"active"`
}

func (p *AddEndpointPayload) Validate() error {
	return validate.Struct(p)
}

type GetEndpointByIDPayload struct {
	ID string `json:"id" validate:"required"`
}

func (p *GetEndpointByIDPayload) Validate() error {
	return validate.Struct(p)
}

type UpdateEndpointPayload struct {
	ID        string `json:"id" validate:"required"`
	TargetURL string `json:"targetUrl" validate:"required,url"`
	Active    bool   `json:"active"`
}

func (p *UpdateEndpointPayload) Validate() error {
	return validate.Struct(p)
}

type DeleteEndpointByIDPayload struct {
	ID string `json:"id" validate:"required"`
}

func (p *DeleteEndpointByIDPayload) Validate() error {
	return validate.Struct(p)
}
