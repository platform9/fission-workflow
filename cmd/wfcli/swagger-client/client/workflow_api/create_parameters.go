package workflow_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/fission/fission-workflows/cmd/wfcli/swagger-client/models"
)

// NewCreateParams creates a new CreateParams object
// with the default values initialized.
func NewCreateParams() *CreateParams {
	var ()
	return &CreateParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateParamsWithTimeout creates a new CreateParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateParamsWithTimeout(timeout time.Duration) *CreateParams {
	var ()
	return &CreateParams{

		timeout: timeout,
	}
}

// NewCreateParamsWithContext creates a new CreateParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateParamsWithContext(ctx context.Context) *CreateParams {
	var ()
	return &CreateParams{

		Context: ctx,
	}
}

/*CreateParams contains all the parameters to send to the API endpoint
for the create operation typically these are written to a http.Request
*/
type CreateParams struct {

	/*Body*/
	Body *models.WorkflowSpec

	timeout time.Duration
	Context context.Context
}

// WithTimeout adds the timeout to the create params
func (o *CreateParams) WithTimeout(timeout time.Duration) *CreateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create params
func (o *CreateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create params
func (o *CreateParams) WithContext(ctx context.Context) *CreateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create params
func (o *CreateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithBody adds the body to the create params
func (o *CreateParams) WithBody(body *models.WorkflowSpec) *CreateParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create params
func (o *CreateParams) SetBody(body *models.WorkflowSpec) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	if o.Body == nil {
		o.Body = new(models.WorkflowSpec)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
