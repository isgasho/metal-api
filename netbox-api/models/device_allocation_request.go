// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DeviceAllocationRequest device allocation request
// swagger:model DeviceAllocationRequest
type DeviceAllocationRequest struct {

	// Additional description for this device in the netbox
	// Min Length: 1
	Description string `json:"description,omitempty"`

	// The desired name for this device in the netbox
	// Required: true
	// Min Length: 1
	Name *string `json:"name"`

	// The operating system name that will be installed on this device
	// Min Length: 1
	Os string `json:"os,omitempty"`

	// The name of the project to assign this device to
	// Required: true
	// Min Length: 1
	Project *string `json:"project"`

	// The name of the tenant to assign this device to
	// Required: true
	// Min Length: 1
	Tenant *string `json:"tenant"`
}

// Validate validates this device allocation request
func (m *DeviceAllocationRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProject(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTenant(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DeviceAllocationRequest) validateDescription(formats strfmt.Registry) error {

	if swag.IsZero(m.Description) { // not required
		return nil
	}

	if err := validate.MinLength("description", "body", string(m.Description), 1); err != nil {
		return err
	}

	return nil
}

func (m *DeviceAllocationRequest) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	if err := validate.MinLength("name", "body", string(*m.Name), 1); err != nil {
		return err
	}

	return nil
}

func (m *DeviceAllocationRequest) validateOs(formats strfmt.Registry) error {

	if swag.IsZero(m.Os) { // not required
		return nil
	}

	if err := validate.MinLength("os", "body", string(m.Os), 1); err != nil {
		return err
	}

	return nil
}

func (m *DeviceAllocationRequest) validateProject(formats strfmt.Registry) error {

	if err := validate.Required("project", "body", m.Project); err != nil {
		return err
	}

	if err := validate.MinLength("project", "body", string(*m.Project), 1); err != nil {
		return err
	}

	return nil
}

func (m *DeviceAllocationRequest) validateTenant(formats strfmt.Registry) error {

	if err := validate.Required("tenant", "body", m.Tenant); err != nil {
		return err
	}

	if err := validate.MinLength("tenant", "body", string(*m.Tenant), 1); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DeviceAllocationRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeviceAllocationRequest) UnmarshalBinary(b []byte) error {
	var res DeviceAllocationRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
