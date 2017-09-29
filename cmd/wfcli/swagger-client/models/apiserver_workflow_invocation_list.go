package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

// ApiserverWorkflowInvocationList apiserver workflow invocation list
// swagger:model apiserverWorkflowInvocationList
type ApiserverWorkflowInvocationList struct {

	// invocations
	Invocations []string `json:"invocations,omitempty"`
}

// Validate validates this apiserver workflow invocation list
func (m *ApiserverWorkflowInvocationList) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInvocations(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApiserverWorkflowInvocationList) validateInvocations(formats strfmt.Registry) error {

	if swag.IsZero(m.Invocations) { // not required
		return nil
	}

	return nil
}
