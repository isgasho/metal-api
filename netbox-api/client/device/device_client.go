// Code generated by go-swagger; DO NOT EDIT.

package device

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new device API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for device API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
LibServerAllocateDevice allocates a device

Allocated a registered device for a tenant
*/
func (a *Client) LibServerAllocateDevice(params *LibServerAllocateDeviceParams) (*LibServerAllocateDeviceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewLibServerAllocateDeviceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "lib.server.allocate_device",
		Method:             "POST",
		PathPattern:        "/allocate-device/{uuid}",
		ProducesMediaTypes: []string{"application/json", "application/problem+json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &LibServerAllocateDeviceReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*LibServerAllocateDeviceOK), nil

}

/*
LibServerRegisterDevice registers a new device

Register a new device to the netbox
*/
func (a *Client) LibServerRegisterDevice(params *LibServerRegisterDeviceParams) (*LibServerRegisterDeviceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewLibServerRegisterDeviceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "lib.server.register_device",
		Method:             "POST",
		PathPattern:        "/register-device/{uuid}",
		ProducesMediaTypes: []string{"application/json", "application/problem+json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &LibServerRegisterDeviceReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*LibServerRegisterDeviceOK), nil

}

/*
LibServerReleaseDevice releases a device

Releases an allocated device
*/
func (a *Client) LibServerReleaseDevice(params *LibServerReleaseDeviceParams) (*LibServerReleaseDeviceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewLibServerReleaseDeviceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "lib.server.release_device",
		Method:             "POST",
		PathPattern:        "/release-device/{uuid}",
		ProducesMediaTypes: []string{"application/json", "application/problem+json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &LibServerReleaseDeviceReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*LibServerReleaseDeviceOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
