// Code generated by go-swagger; DO NOT EDIT.

package devices

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "git.f-i-ts.de/cloud-native/maas/metal-api/netbox-api/models"
)

// NewNetboxAPIProxyAPIDeviceRegisterParams creates a new NetboxAPIProxyAPIDeviceRegisterParams object
// with the default values initialized.
func NewNetboxAPIProxyAPIDeviceRegisterParams() *NetboxAPIProxyAPIDeviceRegisterParams {
	var ()
	return &NetboxAPIProxyAPIDeviceRegisterParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewNetboxAPIProxyAPIDeviceRegisterParamsWithTimeout creates a new NetboxAPIProxyAPIDeviceRegisterParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewNetboxAPIProxyAPIDeviceRegisterParamsWithTimeout(timeout time.Duration) *NetboxAPIProxyAPIDeviceRegisterParams {
	var ()
	return &NetboxAPIProxyAPIDeviceRegisterParams{

		timeout: timeout,
	}
}

// NewNetboxAPIProxyAPIDeviceRegisterParamsWithContext creates a new NetboxAPIProxyAPIDeviceRegisterParams object
// with the default values initialized, and the ability to set a context for a request
func NewNetboxAPIProxyAPIDeviceRegisterParamsWithContext(ctx context.Context) *NetboxAPIProxyAPIDeviceRegisterParams {
	var ()
	return &NetboxAPIProxyAPIDeviceRegisterParams{

		Context: ctx,
	}
}

// NewNetboxAPIProxyAPIDeviceRegisterParamsWithHTTPClient creates a new NetboxAPIProxyAPIDeviceRegisterParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewNetboxAPIProxyAPIDeviceRegisterParamsWithHTTPClient(client *http.Client) *NetboxAPIProxyAPIDeviceRegisterParams {
	var ()
	return &NetboxAPIProxyAPIDeviceRegisterParams{
		HTTPClient: client,
	}
}

/*NetboxAPIProxyAPIDeviceRegisterParams contains all the parameters to send to the API endpoint
for the netbox api proxy api device register operation typically these are written to a http.Request
*/
type NetboxAPIProxyAPIDeviceRegisterParams struct {

	/*Request
	  The device registration body

	*/
	Request *models.DeviceRegistrationRequest
	/*UUID
	  The product serial of the device (unique identifier of this device)

	*/
	UUID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) WithTimeout(timeout time.Duration) *NetboxAPIProxyAPIDeviceRegisterParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) WithContext(ctx context.Context) *NetboxAPIProxyAPIDeviceRegisterParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) WithHTTPClient(client *http.Client) *NetboxAPIProxyAPIDeviceRegisterParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRequest adds the request to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) WithRequest(request *models.DeviceRegistrationRequest) *NetboxAPIProxyAPIDeviceRegisterParams {
	o.SetRequest(request)
	return o
}

// SetRequest adds the request to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) SetRequest(request *models.DeviceRegistrationRequest) {
	o.Request = request
}

// WithUUID adds the uuid to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) WithUUID(uuid string) *NetboxAPIProxyAPIDeviceRegisterParams {
	o.SetUUID(uuid)
	return o
}

// SetUUID adds the uuid to the netbox api proxy api device register params
func (o *NetboxAPIProxyAPIDeviceRegisterParams) SetUUID(uuid string) {
	o.UUID = uuid
}

// WriteToRequest writes these params to a swagger request
func (o *NetboxAPIProxyAPIDeviceRegisterParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Request != nil {
		if err := r.SetBodyParam(o.Request); err != nil {
			return err
		}
	}

	// path param uuid
	if err := r.SetPathParam("uuid", o.UUID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
