package v1

import (
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/metal"
)

type PartitionBase struct {
	MgmtServiceAddress         *string `json:"mgmtserviceaddress" description:"the address to the management service of this partition" optional:"true"`
	ProjectNetworkPrefixLength *int    `json:"projectnetworkprefixlength" description:"the length of project networks for this partition, default 22" optional:"true" minimum:"16" maximum:"30"`
}

type PartitionBootConfiguration struct {
	ImageURL    *string `json:"imageurl" modelDescription:"a partition has a distinct location in a data center, individual entities belong to a partition" description:"the url to download the initrd for the boot image" optional:"true"`
	KernelURL   *string `json:"kernelurl" description:"the url to download the kernel for the boot image" optional:"true"`
	CommandLine *string `json:"commandline" description:"the cmdline to the kernel for the boot image" optional:"true"`
}

type PartitionCreateRequest struct {
	Common
	PartitionBase
	PartitionBootConfiguration PartitionBootConfiguration `json:"bootconfig" description:"the boot configuration of this partition"`
}

type PartitionUpdateRequest struct {
	Common
	MgmtServiceAddress         *string                     `json:"mgmtserviceaddress" description:"the address to the management service of this partition" optional:"true"`
	PartitionBootConfiguration *PartitionBootConfiguration `json:"bootconfig" description:"the boot configuration of this partition" optional:"true"`
}

type PartitionListResponse struct {
	Common
	PartitionBase
	PartitionBootConfiguration PartitionBootConfiguration `json:"bootconfig" description:"the boot configuration of this partition"`
}

type PartitionDetailResponse struct {
	PartitionListResponse
	Timestamps
}

func NewPartitionDetailResponse(p *metal.Partition) *PartitionDetailResponse {
	if p == nil {
		return nil
	}
	return &PartitionDetailResponse{
		PartitionListResponse: *NewPartitionListResponse(p),
		Timestamps: Timestamps{
			Created: p.Created,
			Changed: p.Changed,
		},
	}
}

func NewPartitionListResponse(p *metal.Partition) *PartitionListResponse {
	if p == nil {
		return nil
	}
	return &PartitionListResponse{
		Common: Common{
			Identifiable: Identifiable{
				ID: p.ID,
			},
			Describeable: Describeable{
				Name:        &p.Name,
				Description: &p.Description,
			},
		},
		PartitionBase: PartitionBase{
			MgmtServiceAddress:         &p.MgmtServiceAddress,
			ProjectNetworkPrefixLength: &p.ProjectNetworkPrefixLength,
		},
		PartitionBootConfiguration: PartitionBootConfiguration{
			ImageURL:    &p.BootConfiguration.ImageURL,
			KernelURL:   &p.BootConfiguration.KernelURL,
			CommandLine: &p.BootConfiguration.CommandLine,
		},
	}
}
