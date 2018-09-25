package maas

import "time"

type Device struct {
	ID          string    `json:"id" description:"a unique ID" unique:"true" readOnly:"true" modelDescription:"A device representing a bare metal machine."`
	Name        string    `json:"name" description:"the name of the device"`
	Description string    `json:"description,omitempty" description:"a description for this machine" optional:"true"`
	Created     time.Time `json:"created" description:"the creation time of this machine" optional:"true" readOnly:"true"`
	Changed     time.Time `json:"changed" description:"the last changed timestamp" optional:"true" readOnly:"true"`
	Project     string    `json:"project" description:"the project that this device is assigned to"`
	Facility    Facility  `json:"facility" description:"the facility assigned to this device"`
	Image       Image     `json:"image" description:"the image assigned to this device"`
	Size        Size      `json:"size" description:"the size of this device"`
}
