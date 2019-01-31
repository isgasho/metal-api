package metal

// A MacAddress is the type for mac adresses. When using a
// custom type, we cannot use strings directly.
type MacAddress string

// Nic information.
type Nic struct {
	MacAddress MacAddress `json:"mac"  description:"the mac address of this network interface" rethinkdb:"macAddress"`
	Name       string     `json:"name"  description:"the name of this network interface" rethinkdb:"name"`
	Vrf        string     `json:"vrf" description:"the vrf this network interface is part of" rethinkdb:"vrf" optional:"true"`
	Neighbors  Nics       `json:"neighbors" description:"the neighbors visible to this network interface" rethinkdb:"neighbors"`
}

type Vrf struct {
	ID     uint
	Tenant string
}

// Nics is a list of nics.
type Nics []Nic

// ByMac creates a indexed map from a nic list.
func (nics Nics) ByMac() map[MacAddress]*Nic {
	res := make(map[MacAddress]*Nic)
	for i, n := range nics {
		res[n.MacAddress] = &nics[i]
	}
	return res
}
