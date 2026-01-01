package endpoint

import (
	"github.com/go-playground/validator/v10"
)

// CRUD (Create, read, update and delete)

// CREATE ---------------------------------------------------

type AddEndpointPayload struct {
	ID        string `json:"id" validate:"required"`
	TargetUrl string `json:"targetUrl" validate:"required, http_url"`
	Active    bool   `json:"active" validate:"required"`
}

func (p *AddEndpointPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------

// READ ---------------------------------------------------

type GetEndpointbyIDPayload struct {
	ID string `json:"id" validate:"required"`
}

func (p *GetEndpointbyIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------

// UPDATE ---------------------------------------------------

type UpdateEndpointPayload struct {
	ID        string `json:"id" validate:"required"`
	TargetUrl string `json:"targetUrl" validate:"required, http_url"`
	Active    bool   `json:"active" validate:"required"`
}

func (p *UpdateEndpointPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------

// DELETE---------------------------------------------------

type DeleteEndpointbyIDPayload struct {
	ID string `json:"id" validate:"required"`
}

func (p *DeleteEndpointbyIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
