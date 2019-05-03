// Code generated by go-swagger; DO NOT EDIT.

package ips

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "git.f-i-ts.de/cloud-native/metal/metal-api/netbox-api/models"
)

// NetboxAPIProxyAPIDeleteIPReader is a Reader for the NetboxAPIProxyAPIDeleteIP structure.
type NetboxAPIProxyAPIDeleteIPReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *NetboxAPIProxyAPIDeleteIPReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewNetboxAPIProxyAPIDeleteIPOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewNetboxAPIProxyAPIDeleteIPDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewNetboxAPIProxyAPIDeleteIPOK creates a NetboxAPIProxyAPIDeleteIPOK with default headers values
func NewNetboxAPIProxyAPIDeleteIPOK() *NetboxAPIProxyAPIDeleteIPOK {
	return &NetboxAPIProxyAPIDeleteIPOK{}
}

/*NetboxAPIProxyAPIDeleteIPOK handles this case with default header values.

OK
*/
type NetboxAPIProxyAPIDeleteIPOK struct {
}

func (o *NetboxAPIProxyAPIDeleteIPOK) Error() string {
	return fmt.Sprintf("[DELETE /ips][%d] netboxApiProxyApiDeleteIpOK ", 200)
}

func (o *NetboxAPIProxyAPIDeleteIPOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewNetboxAPIProxyAPIDeleteIPDefault creates a NetboxAPIProxyAPIDeleteIPDefault with default headers values
func NewNetboxAPIProxyAPIDeleteIPDefault(code int) *NetboxAPIProxyAPIDeleteIPDefault {
	return &NetboxAPIProxyAPIDeleteIPDefault{
		_statusCode: code,
	}
}

/*NetboxAPIProxyAPIDeleteIPDefault handles this case with default header values.

Problem
*/
type NetboxAPIProxyAPIDeleteIPDefault struct {
	_statusCode int

	Payload *models.Problem
}

// Code gets the status code for the netbox api proxy api delete ip default response
func (o *NetboxAPIProxyAPIDeleteIPDefault) Code() int {
	return o._statusCode
}

func (o *NetboxAPIProxyAPIDeleteIPDefault) Error() string {
	return fmt.Sprintf("[DELETE /ips][%d] netbox_api_proxy.api.delete_ip default  %+v", o._statusCode, o.Payload)
}

func (o *NetboxAPIProxyAPIDeleteIPDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Problem)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
