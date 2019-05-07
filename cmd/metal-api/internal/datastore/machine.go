package datastore

import (
	"fmt"

	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/metal"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// Some predefined error values.
var (
	ErrNoMachineAvailable = fmt.Errorf("no machine available")
)

// FindMachine returns the machine with the given ID. If there is no
// such machine a metal.NotFound will be returned.
func (rs *RethinkStore) FindMachine(id string) (*metal.Machine, error) {
	var m metal.Machine
	err := rs.findEntityByID(rs.machineTable(), &m, id)
	return &m, err
}

// FindMachineAllowNil returns the machine with the given ID. If there is no
// such machine nil will be returned.
func (rs *RethinkStore) FindMachineAllowNil(id string) (*metal.Machine, error) {
	var m metal.Machine
	err := rs.findEntityByIDAllowNil(rs.machineTable(), &m, id)
	if m.ID != "" {
		return &m, err
	}
	return nil, err
}

// ListMachines returns all machines.
func (rs *RethinkStore) ListMachines() ([]metal.Machine, error) {
	ms := make([]metal.Machine, 0)
	err := rs.listEntities(rs.machineTable(), &ms)
	return ms, err
}

// SearchMachine returns the machines filtered by the given search filter.
func (rs *RethinkStore) SearchMachine(mac string) ([]metal.Machine, error) {
	searchFilter := func(row r.Term) r.Term {
		return row
	}

	if mac != "" {
		searchFilter = func(row r.Term) r.Term {
			return row.Field("hardware").Field("network_interfaces").Map(func(nic r.Term) r.Term {
				return nic.Field("macAddress")
			}).Contains(r.Expr(mac))
		}
	}

	var ms []metal.Machine
	err := rs.searchEntities(rs.machineTable(), searchFilter, &ms)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

// CreateMachine creates a new machine in the database as "unallocated new machines".
// If the given machine has an allocation, the function returns an error because
// allocated machines cannot be created. If there is already a machine with the
// given ID in the database it will be replaced the the given machine.
// CreateNetwork creates a new network.
func (rs *RethinkStore) CreateMachine(m *metal.Machine) error {
	if m.Allocation != nil {
		return fmt.Errorf("a machine cannot be created when it is allocated: %q: %+v", m.ID, *m.Allocation)
	}
	return rs.createEntity(rs.machineTable(), m)
}

// DeleteMachine removes a machine from the database.
func (rs *RethinkStore) DeleteMachine(m *metal.Machine) error {
	return rs.deleteEntityByID(rs.machineTable(), m.GetID())
}

// UpdateMachine replaces a machine in the database if the 'changed' field of
// the old value equals the 'changed' field of the recored in the database.
func (rs *RethinkStore) UpdateMachine(oldMachine *metal.Machine, newMachine *metal.Machine) error {
	return rs.updateEntity(rs.machineTable(), newMachine, oldMachine)
}

// InsertWaitingMachine adds a machine to the wait table.
func (rs *RethinkStore) InsertWaitingMachine(m *metal.Machine) error {
	// does not prohibit concurrent wait calls for the same UUID
	_, err := rs.waitTable().Insert(m, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(rs.session)
	if err != nil {
		return fmt.Errorf("cannot insert machine into wait table: %v", err)
	}
	return nil
}

// ListenWaitTable listens on changes on the wait table for a given machine id.
func (rs *RethinkStore) WaitForMachineAllocation(id string) (*metal.Machine, error) {
	ch, err := rs.waitTable().Get(id).Changes().Run(rs.session)
	if err != nil {
		rs.SugaredLogger.Errorw("cannot wait for allocation", "error", err)
		// simply return so this machine will not be allocated
		// the normal timeout-behaviour of the allocator will
		// occur without an allocation
		return nil, err
	}
	type responseType struct {
		NewVal metal.Machine `rethinkdb:"new_val"`
		OldVal metal.Machine `rethinkdb:"old_val"`
	}
	var response responseType
	for ch.Next(&response) {
		rs.SugaredLogger.Infow("machine changed", "response", response)
		if response.NewVal.ID == "" {
			// the entry was deleted, no wait any more
			return nil, fmt.Errorf("machine %s not available any more", id)
		}
		return &response.NewVal, nil
	}
	return nil, fmt.Errorf("no machine found even though change event was received")
}

// RemoveWaitingMachine removes a machine from the wait table.
func (rs *RethinkStore) RemoveWaitingMachine(id string) error {
	_, err := rs.waitTable().Get(id).Delete().RunWrite(rs.session)
	return err
}

// UpdateWaitingMachine updates a machine in the wait table with the given machine
func (rs *RethinkStore) UpdateWaitingMachine(m *metal.Machine) error {
	_, err := rs.waitTable().Get(m.ID).Update(m).RunWrite(rs.session)
	return err
}

// FindAvailableMachine returns an available machine from the wait table.
func (rs *RethinkStore) FindAvailableMachine(partitionid, sizeid string) (*metal.Machine, error) {
	query := map[string]interface{}{
		"allocation":  nil,
		"partitionid": partitionid,
		"sizeid":      sizeid,
		"state": map[string]interface{}{
			"value": "",
		},
	}
	var available []metal.Machine
	err := rs.searchEntities(rs.waitTable(), query, &available)
	if err != nil {
		return nil, fmt.Errorf("cannot find free machine: %v", err)
	}

	if len(available) < 1 {
		return nil, ErrNoMachineAvailable
	}

	return &available[0], nil
}
