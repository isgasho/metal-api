// Code generated by go-swagger; DO NOT EDIT.

package networks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "git.f-i-ts.de/cloud-native/metal/metal-api/netbox-api/models"
)

// NetboxAPIProxyAPIAllocateNetworkReader is a Reader for the NetboxAPIProxyAPIAllocateNetwork structure.
type NetboxAPIProxyAPIAllocateNetworkReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *NetboxAPIProxyAPIAllocateNetworkReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewNetboxAPIProxyAPIAllocateNetworkOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewNetboxAPIProxyAPIAllocateNetworkDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewNetboxAPIProxyAPIAllocateNetworkOK creates a NetboxAPIProxyAPIAllocateNetworkOK with default headers values
func NewNetboxAPIProxyAPIAllocateNetworkOK() *NetboxAPIProxyAPIAllocateNetworkOK {
	return &NetboxAPIProxyAPIAllocateNetworkOK{}
}

/*NetboxAPIProxyAPIAllocateNetworkOK handles this case with default header values.

OK
*/
type NetboxAPIProxyAPIAllocateNetworkOK struct {
	Payload *models.CIDR
}

func (o *NetboxAPIProxyAPIAllocateNetworkOK) Error() string {
	return fmt.Sprintf("[POST /networks/allocate][%d] netboxApiProxyApiAllocateNetworkOK  %+v", 200, o.Payload)
}

func (o *NetboxAPIProxyAPIAllocateNetworkOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CIDR)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewNetboxAPIProxyAPIAllocateNetworkDefault creates a NetboxAPIProxyAPIAllocateNetworkDefault with default headers values
func NewNetboxAPIProxyAPIAllocateNetworkDefault(code int) *NetboxAPIProxyAPIAllocateNetworkDefault {
	return &NetboxAPIProxyAPIAllocateNetworkDefault{
		_statusCode: code,
	}
}

/*NetboxAPIProxyAPIAllocateNetworkDefault handles this case with default header values.

Problem
*/
type NetboxAPIProxyAPIAllocateNetworkDefault struct {
	_statusCode int

	Payload *models.Problem
}

// Code gets the status code for the netbox api proxy api allocate network default response
func (o *NetboxAPIProxyAPIAllocateNetworkDefault) Code() int {
	return o._statusCode
}

func (o *NetboxAPIProxyAPIAllocateNetworkDefault) Error() string {
	return fmt.Sprintf("[POST /networks/allocate][%d] netbox_api_proxy.api.allocate_network default  %+v", o._statusCode, o.Payload)
}

func (o *NetboxAPIProxyAPIAllocateNetworkDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Problem)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}