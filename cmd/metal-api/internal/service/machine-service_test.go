package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/datastore"
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/ipam"
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/metal"
	v1 "git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/service/v1"
	"git.f-i-ts.de/cloud-native/metal/metal-api/cmd/metal-api/internal/testdata"
	"git.f-i-ts.de/cloud-native/metallib/httperrors"
	"github.com/emicklei/go-restful"
	goipam "github.com/metal-pod/go-ipam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type emptyPublisher struct {
	doPublish func(topic string, data interface{}) error
}

func (p *emptyPublisher) Publish(topic string, data interface{}) error {
	if p.doPublish != nil {
		return p.doPublish(topic, data)
	}
	return nil
}

func (p *emptyPublisher) CreateTopic(topic string) error {
	return nil
}

func TestGetMachines(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
	container := restful.NewContainer().Add(machineservice)
	req := httptest.NewRequest("GET", "/v1/machine", nil)
	container = injectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result []v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Len(t, result, len(testdata.TestMachines))
	require.Equal(t, testdata.M1.ID, result[0].ID)
	require.Equal(t, testdata.M1.Allocation.Name, result[0].Allocation.Name)
	require.Equal(t, testdata.Sz1.Name, *result[0].Size.Name)
	require.Equal(t, testdata.Partition1.Name, *result[0].Partition.Name)
	require.Equal(t, testdata.M2.ID, result[1].ID)
}

func TestRegisterMachine(t *testing.T) {
	data := []struct {
		name                 string
		uuid                 string
		partitionid          string
		numcores             int
		memory               int
		dbpartitions         []metal.Partition
		dbsizes              []metal.Size
		dbmachines           []metal.Machine
		expectedStatus       int
		expectedErrorMessage string
		expectedSizeName     string
	}{
		{
			name:             "insert new",
			uuid:             "0",
			partitionid:      "0",
			dbpartitions:     []metal.Partition{testdata.Partition1},
			dbsizes:          []metal.Size{testdata.Sz1},
			numcores:         1,
			memory:           100,
			expectedStatus:   http.StatusCreated,
			expectedSizeName: testdata.Sz1.Name,
		},
		{
			name:             "insert existing",
			uuid:             "1",
			partitionid:      "1",
			dbpartitions:     []metal.Partition{testdata.Partition1},
			dbsizes:          []metal.Size{testdata.Sz1},
			dbmachines:       []metal.Machine{testdata.M1},
			numcores:         1,
			memory:           100,
			expectedStatus:   http.StatusOK,
			expectedSizeName: testdata.Sz1.Name,
		},
		{
			name:                 "empty uuid",
			uuid:                 "",
			partitionid:          "1",
			dbpartitions:         []metal.Partition{testdata.Partition1},
			dbsizes:              []metal.Size{testdata.Sz1},
			expectedStatus:       http.StatusUnprocessableEntity,
			expectedErrorMessage: "uuid cannot be empty",
		},
		{
			name:                 "empty partition",
			uuid:                 "1",
			partitionid:          "",
			dbpartitions:         nil,
			dbsizes:              []metal.Size{testdata.Sz1},
			expectedStatus:       http.StatusNotFound,
			expectedErrorMessage: "no partition with id \"\" found",
		},
		{
			name:             "new with unknown size",
			uuid:             "0",
			partitionid:      "1",
			dbpartitions:     []metal.Partition{testdata.Partition1},
			dbsizes:          []metal.Size{testdata.Sz1},
			numcores:         2,
			memory:           100,
			expectedStatus:   http.StatusCreated,
			expectedSizeName: metal.UnknownSize.Name,
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			ds, mock := datastore.InitMockDB()
			mock.On(r.DB("mockdb").Table("partition").Get(test.partitionid)).Return(test.dbpartitions, nil)

			if len(test.dbmachines) > 0 {
				mock.On(r.DB("mockdb").Table("size").Get(test.dbmachines[0].SizeID)).Return([]metal.Size{testdata.Sz1}, nil)
				mock.On(r.DB("mockdb").Table("machine").Get(test.dbmachines[0].ID).Replace(r.MockAnything())).Return(testdata.EmptyResult, nil)
			} else {
				mock.On(r.DB("mockdb").Table("machine").Get("0")).Return(nil, nil)
				mock.On(r.DB("mockdb").Table("machine").Insert(r.MockAnything(), r.InsertOpts{
					Conflict: "replace",
				})).Return(testdata.EmptyResult, nil)
			}
			mock.On(r.DB("mockdb").Table("size").Get(metal.UnknownSize.ID)).Return([]metal.Size{*metal.UnknownSize}, nil)
			mock.On(r.DB("mockdb").Table("switch").Filter(r.MockAnything(), r.FilterOpts{})).Return([]metal.Switch{}, nil)
			mock.On(r.DB("mockdb").Table("event").Filter(r.MockAnything(), r.FilterOpts{})).Return([]metal.ProvisioningEventContainer{}, nil)
			mock.On(r.DB("mockdb").Table("event").Insert(r.MockAnything(), r.InsertOpts{})).Return(testdata.EmptyResult, nil)
			testdata.InitMockDBData(mock)

			registerRequest := &v1.MachineRegisterRequest{
				UUID:        test.uuid,
				PartitionID: test.partitionid,
				RackID:      "1",
				IPMI: v1.MachineIPMI{
					Address:    testdata.IPMI1.Address,
					Interface:  testdata.IPMI1.Interface,
					MacAddress: testdata.IPMI1.MacAddress,
					Fru: v1.MachineFru{
						ChassisPartNumber:   &testdata.IPMI1.Fru.ChassisPartNumber,
						ChassisPartSerial:   &testdata.IPMI1.Fru.ChassisPartSerial,
						BoardMfg:            &testdata.IPMI1.Fru.BoardMfg,
						BoardMfgSerial:      &testdata.IPMI1.Fru.BoardMfgSerial,
						BoardPartNumber:     &testdata.IPMI1.Fru.BoardPartNumber,
						ProductManufacturer: &testdata.IPMI1.Fru.ProductManufacturer,
						ProductPartNumber:   &testdata.IPMI1.Fru.ProductPartNumber,
						ProductSerial:       &testdata.IPMI1.Fru.ProductSerial,
					},
				},
				Hardware: v1.MachineHardwareExtended{
					MachineHardwareBase: v1.MachineHardwareBase{
						CPUCores: test.numcores,
						Memory:   uint64(test.memory),
					},
				},
			}

			js, _ := json.Marshal(registerRequest)
			body := bytes.NewBuffer(js)
			machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
			container := restful.NewContainer().Add(machineservice)
			req := httptest.NewRequest("POST", "/v1/machine/register", body)
			req.Header.Add("Content-Type", "application/json")
			container = injectEditor(container, req)
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, test.expectedStatus, resp.StatusCode, w.Body.String())

			if test.expectedStatus > 300 {
				var result httperrors.HTTPErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Regexp(t, test.expectedErrorMessage, result.Message)
			} else {
				var result v1.MachineResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				expectedid := "0"
				if len(test.dbmachines) > 0 {
					expectedid = test.dbmachines[0].ID
				}
				require.Equal(t, expectedid, result.ID)
				require.Equal(t, "1", result.RackID)
				require.Equal(t, test.expectedSizeName, *result.Size.Name)
				require.Equal(t, testdata.Partition1.Name, *result.Partition.Name)
			}
		})
	}
}

func TestMachineIPMI(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	data := []struct {
		name           string
		machine        *metal.Machine
		wantStatusCode int
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:           "retrieve machine1 ipmi",
			machine:        &testdata.M1,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "retrieve machine2 ipmi",
			machine:        &testdata.M2,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "retrieve unknown machine ipmi",
			machine:        &metal.Machine{Base: metal.Base{ID: "999"}},
			wantStatusCode: http.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {

			machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
			container := restful.NewContainer().Add(machineservice)

			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/machine/%s/ipmi", test.machine.ID), nil)
			req.Header.Add("Content-Type", "application/json")
			container = injectViewer(container, req)
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, test.wantStatusCode, resp.StatusCode, w.Body.String())

			if test.wantErr {
				var result httperrors.HTTPErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Equal(t, test.wantStatusCode, result.StatusCode)
				if test.wantErrMessage != "" {
					require.Regexp(t, test.wantErrMessage, result.Message)
				}
			} else {
				var result v1.MachineIPMI
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Equal(t, test.machine.IPMI.Address, result.Address)
				require.Equal(t, test.machine.IPMI.Interface, result.Interface)
				require.Equal(t, test.machine.IPMI.User, result.User)
				require.Equal(t, test.machine.IPMI.Password, result.Password)
				require.Equal(t, test.machine.IPMI.MacAddress, result.MacAddress)

				require.Equal(t, test.machine.IPMI.Fru.ChassisPartNumber, *result.Fru.ChassisPartNumber)
				require.Equal(t, test.machine.IPMI.Fru.ChassisPartSerial, *result.Fru.ChassisPartSerial)
				require.Equal(t, test.machine.IPMI.Fru.BoardMfg, *result.Fru.BoardMfg)
				require.Equal(t, test.machine.IPMI.Fru.BoardMfgSerial, *result.Fru.BoardMfgSerial)
				require.Equal(t, test.machine.IPMI.Fru.BoardPartNumber, *result.Fru.BoardPartNumber)
				require.Equal(t, test.machine.IPMI.Fru.ProductManufacturer, *result.Fru.ProductManufacturer)
				require.Equal(t, test.machine.IPMI.Fru.ProductPartNumber, *result.Fru.ProductPartNumber)
				require.Equal(t, test.machine.IPMI.Fru.ProductSerial, *result.Fru.ProductSerial)
			}
		})
	}
}

func TestFinalizeMachineAllocation(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	data := []struct {
		name           string
		machineID      string
		wantStatusCode int
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:           "finalize successfully",
			machineID:      "1",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "finalize unknown machine",
			machineID:      "999",
			wantStatusCode: http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:           "finalize unallocated machine",
			machineID:      "3",
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErr:        true,
			wantErrMessage: "the machine \"3\" is not allocated",
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {

			machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
			container := restful.NewContainer().Add(machineservice)

			finalizeRequest := v1.MachineFinalizeAllocationRequest{
				ConsolePassword: "blubber",
			}

			js, _ := json.Marshal(finalizeRequest)
			body := bytes.NewBuffer(js)
			req := httptest.NewRequest("POST", fmt.Sprintf("/v1/machine/%s/finalize-allocation", test.machineID), body)
			req.Header.Add("Content-Type", "application/json")
			container = injectEditor(container, req)
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, test.wantStatusCode, resp.StatusCode, w.Body.String())

			if test.wantErr {
				var result httperrors.HTTPErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Equal(t, test.wantStatusCode, result.StatusCode)
				if test.wantErrMessage != "" {
					require.Regexp(t, test.wantErrMessage, result.Message)
				}
			} else {
				var result v1.MachineResponse
				err := json.NewDecoder(resp.Body).Decode(&result)

				require.Nil(t, err)
				require.Equal(t, finalizeRequest.ConsolePassword, *result.Allocation.ConsolePassword)
			}
		})
	}
}
func TestSetMachineState(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
	container := restful.NewContainer().Add(machineservice)

	stateRequest := v1.MachineState{
		Value:       string(metal.ReservedState),
		Description: "blubber",
	}
	js, _ := json.Marshal(stateRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/machine/1/state", body)
	req.Header.Add("Content-Type", "application/json")
	container = injectEditor(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, "1", result.ID)
	require.Equal(t, string(metal.ReservedState), result.State.Value)
	require.Equal(t, "blubber", result.State.Description)

}

func TestGetMachine(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
	container := restful.NewContainer().Add(machineservice)
	req := httptest.NewRequest("GET", "/v1/machine/1", nil)
	container = injectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.M1.ID, result.ID)
	require.Equal(t, testdata.M1.Allocation.Name, result.Allocation.Name)
	require.Equal(t, testdata.Sz1.Name, *result.Size.Name)
	require.Equal(t, testdata.Img1.Name, *result.Allocation.Image.Name)
	require.Equal(t, testdata.Partition1.Name, *result.Partition.Name)
}

func TestGetMachineNotFound(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
	container := restful.NewContainer().Add(machineservice)
	req := httptest.NewRequest("GET", "/v1/machine/999", nil)
	container = injectEditor(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusNotFound, resp.StatusCode, w.Body.String())
}

func TestFreeMachine(t *testing.T) {
	// TODO: Add tests for IPAM, verifying that networks are cleaned up properly

	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	pub := &emptyPublisher{}
	events := []string{"1-machine", "1-switch"}
	eventidx := 0
	pub.doPublish = func(topic string, data interface{}) error {
		require.Equal(t, events[eventidx], topic)
		eventidx++
		if eventidx == 0 {
			dv := data.(metal.MachineEvent)
			require.Equal(t, "1", dv.Old.ID)
		}
		return nil
	}

	machineservice := NewMachine(ds, pub, ipam.New(goipam.New()))
	container := restful.NewContainer().Add(machineservice)
	req := httptest.NewRequest("DELETE", "/v1/machine/1/free", nil)
	container = injectEditor(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.M1.ID, result.ID)
	require.Nil(t, result.Allocation)
}

func TestSearchMachine(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	mock.On(r.DB("mockdb").Table("machine").Filter(r.MockAnything())).Return([]interface{}{testdata.M1}, nil)
	testdata.InitMockDBData(mock)

	machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
	container := restful.NewContainer().Add(machineservice)
	requestJSON := fmt.Sprintf("{%q:[%q]}", "nics_mac_addresses", "1")
	req := httptest.NewRequest("POST", "/v1/machine/find", bytes.NewBufferString(requestJSON))
	req.Header.Add("Content-Type", "application/json")
	container = injectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var results []v1.MachineResponse
	err := json.NewDecoder(resp.Body).Decode(&results)

	require.Nil(t, err)
	require.Len(t, results, 1)
	result := results[0]
	require.Equal(t, testdata.M1.ID, result.ID)
	require.Equal(t, testdata.M1.Allocation.Name, result.Allocation.Name)
	require.Equal(t, testdata.Sz1.Name, *result.Size.Name)
	require.Equal(t, testdata.Img1.Name, *result.Allocation.Image.Name)
	require.Equal(t, testdata.Partition1.Name, *result.Partition.Name)
}

func TestAddProvisioningEvent(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	machineservice := NewMachine(ds, &emptyPublisher{}, ipam.New(goipam.New()))
	container := restful.NewContainer().Add(machineservice)
	event := &metal.ProvisioningEvent{
		Event:   metal.ProvisioningEventPreparing,
		Message: "starting metal-hammer",
	}
	js, _ := json.Marshal(event)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/machine/1/event", body)
	container = injectEditor(container, req)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.MachineRecentProvisioningEvents
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, "0", result.IncompleteProvisioningCycles)
	require.Len(t, result.Events, 1)
	if len(result.Events) > 0 {
		require.Equal(t, "starting metal-hammer", result.Events[0].Message)
		require.Equal(t, string(metal.ProvisioningEventPreparing), result.Events[0].Event)
	}
}

func TestOnMachine(t *testing.T) {

	data := []struct {
		cmd      metal.MachineCommand
		endpoint string
		param    string
	}{
		{
			cmd:      metal.MachineOnCmd,
			endpoint: "on",
		},
		{
			cmd:      metal.MachineOffCmd,
			endpoint: "off",
		},
		{
			cmd:      metal.MachineResetCmd,
			endpoint: "reset",
		},
		{
			cmd:      metal.MachineBiosCmd,
			endpoint: "bios",
		},
		{
			cmd:      metal.MachineLedOnCmd,
			endpoint: "led-on",
		},
		{
			cmd:      metal.MachineLedOffCmd,
			endpoint: "led-off",
		},
	}

	for _, d := range data {
		t.Run("cmd_"+d.endpoint, func(t *testing.T) {
			ds, mock := datastore.InitMockDB()
			testdata.InitMockDBData(mock)

			pub := &emptyPublisher{}
			pub.doPublish = func(topic string, data interface{}) error {
				require.Equal(t, "1-machine", topic)
				dv := data.(metal.MachineEvent)
				require.Equal(t, d.cmd, dv.Cmd.Command)
				require.Equal(t, "1", dv.Cmd.Target.ID)
				return nil
			}

			machineservice := NewMachine(ds, pub, ipam.New(goipam.New()))

			js, _ := json.Marshal([]string{d.param})
			body := bytes.NewBuffer(js)
			container := restful.NewContainer().Add(machineservice)
			req := httptest.NewRequest("POST", "/v1/machine/1/power/"+d.endpoint, body)
			container = injectEditor(container, req)
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			container.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
		})
	}
}

func TestParsePublicKey(t *testing.T) {
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDi4+MA0u/luzH2iaKnBTHzo+BEmV1MsdWtPtAps9ccD1vF94AqKtV6mm387ZhamfWUfD1b3Q5ftk56ekwZgHbk6PIUb/W4GrBD4uslTL2lzNX9v0Njo9DfapDKv4Tth6Qz5ldUb6z7IuyDmWqn3FbIPo4LOZxJ9z/HUWyau8+JMSpwIyzp2S0Gtm/pRXhbkZlr4h9jGApDQICPFGBWFEVpyOOjrS8JnEC8YzUszvbj5W1CH6Sn/DtxW0/CTAWwcjIAYYV8GlouWjjALqmjvpxO3F5kvQ1xR8IYrD86+cSCQSP4TpehftzaQzpY98fcog2YkEra+1GCY456cVSUhe1X"
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKey))
	require.Nil(t, err)

	pubKey = ""
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(pubKey))
	require.NotNil(t, err)

	pubKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDi4+MA0u/luzH2iaKnBTHzo+BEmV1MsdWtPtAps9ccD1vF94AqKtV6mm387ZhamfWUfD1b3Q5ftk56ekwZgHbk6PIUb/W4GrBD4uslTL2lzNX9v0Njo9DfapDKv4Tth6Qz5ldUb6z7IuyDmWqn3FbIPo4LOZxJ9z/HUWyau8+JMSpwIyzp2S0Gtm/pRXhbkZlr4h9jGApDQICPFGBWFEVpyOOjrS8JnEC8YzUszvbj5W1CH6Sn/DtxW0/CTAWwcjIAYYV8GlouWjjALqmjvpxO3F5kvQ1xR8IYrD86+cSCQSP4TpehftzaQzpY98fcog2YkEra+1GCY456cVSUhe1"
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(pubKey))
	require.NotNil(t, err)

	pubKey = "AAAAB3NzaC1yc2EAAAADAQABAAABAQDi4+MA0u/luzH2iaKnBTHzo+BEmV1MsdWtPtAps9ccD1vF94AqKtV6mm387ZhamfWUfD1b3Q5ftk56ekwZgHbk6PIUb/W4GrBD4uslTL2lzNX9v0Njo9DfapDKv4Tth6Qz5ldUb6z7IuyDmWqn3FbIPo4LOZxJ9z/HUWyau8+JMSpwIyzp2S0Gtm/pRXhbkZlr4h9jGApDQICPFGBWFEVpyOOjrS8JnEC8YzUszvbj5W1CH6Sn/DtxW0/CTAWwcjIAYYV8GlouWjjALqmjvpxO3F5kvQ1xR8IYrD86+cSCQSP4TpehftzaQzpY98fcog2YkEra+1GCY456cVSUhe1X"
	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(pubKey))
	require.NotNil(t, err)
}

func Test_validateAllocationSpec(t *testing.T) {
	assert := assert.New(t)
	trueValue := true
	falseValue := false

	tests := []struct {
		spec     machineAllocationSpec
		isError  bool
		name     string
		expected string
	}{
		{
			spec: machineAllocationSpec{
				Tenant:     "gopher",
				UUID:       "gopher-uuid",
				IsFirewall: false,
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID: "network",
					},
				},
				IPs: []string{"1.2.3.4"},
			},
			isError:  false,
			expected: "",
			name:     "auto acquire network and additional ip",
		},
		{
			spec: machineAllocationSpec{
				Tenant: "gopher",
				UUID:   "gopher-uuid",
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     "network",
						AutoAcquireIP: &trueValue},
				},
			},
			isError: false,
			name:    "good case (explicit network)",
		},
		{
			spec: machineAllocationSpec{
				Tenant:     "gopher",
				UUID:       "gopher-uuid",
				IsFirewall: false,
			},
			isError:  false,
			expected: "",
			name:     "good case (no network)",
		},
		{
			spec: machineAllocationSpec{
				Tenant:      "gopher",
				PartitionID: "42",
				SizeID:      "42",
			},
			isError: false,
			name:    "partition and size id for absent uuid",
		},
		{
			spec: machineAllocationSpec{
				Tenant:      "gopher",
				PartitionID: "42",
			},
			isError:  true,
			expected: "when no machine id is given, a size id must be specified",
			name:     "missing size id",
		},
		{
			spec: machineAllocationSpec{
				Tenant: "gopher",
				SizeID: "42",
			},
			isError:  true,
			expected: "when no machine id is given, a partition id must be specified",
			name:     "missing partition id",
		},
		{
			spec: machineAllocationSpec{
				UUID: "42",
			},
			isError:  true,
			expected: "no tenant given",
			name:     "absent tenant",
		},
		{
			spec: machineAllocationSpec{
				Tenant: "gopher",
			},
			isError:  true,
			expected: "when no machine id is given, a partition id must be specified",
			name:     "absent uuid",
		},
		{
			spec: machineAllocationSpec{
				Tenant:     "gopher",
				UUID:       "gopher-uuid",
				IsFirewall: false,
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID:     "network",
						AutoAcquireIP: &falseValue,
					},
				},
			},
			isError:  true,
			expected: "missing ip(s) for network(s) without automatic ip allocation",
			name:     "missing ip definition for noauto network",
		},
		{
			spec: machineAllocationSpec{
				Tenant: "gopher",
				UUID:   "42",
				IPs:    []string{"42"},
			},
			isError:  true,
			expected: `"42" is not a valid IP address`,
			name:     "illegal ip",
		},
		{
			spec: machineAllocationSpec{
				Tenant:     "gopher",
				UUID:       "42",
				IsFirewall: true,
			},
			isError:  true,
			expected: "when no ip is given at least one auto acquire network must be specified",
			name:     "missing network/ ip in case of firewall",
		},
		{
			spec: machineAllocationSpec{
				Tenant:     "gopher",
				UUID:       "42",
				SSHPubKeys: []string{"42"},
			},
			isError:  true,
			expected: `invalid public SSH key: 42`,
			name:     "invalid ssh",
		},
		{
			spec: machineAllocationSpec{
				Tenant:     "gopher",
				UUID:       "gopher-uuid",
				IsFirewall: false,
				Networks: []v1.MachineAllocationNetwork{
					{
						NetworkID: "network",
					},
				},
			},
			isError:  false,
			expected: "",
			name:     "implicit auto acquire network",
		},
	}

	for _, test := range tests {
		err := validateAllocationSpec(&test.spec)
		if test.isError {
			assert.Error(err, "Test: %s", test.name)
			assert.EqualError(err, test.expected, "Test: %s", test.name)
		} else {
			assert.NoError(err, "Test: %s", test.name)
		}
	}

}

func Test_additionalTags(t *testing.T) {

	//
	networks := []*metal.MachineNetwork{}
	network := &metal.MachineNetwork{
		Primary:  true,
		IPs:      []string{"1.2.3.4"},
		Prefixes: []string{"1.2.0.0/22", "1.2.2.0/22", "2.3.4.0/24"},
		ASN:      1203874,
	}
	networks = append(networks, network)

	networks26 := []*metal.MachineNetwork{}
	network26 := &metal.MachineNetwork{
		Primary:  true,
		IPs:      []string{"1.2.2.67"},
		Prefixes: []string{"1.2.1.0/28", "1.2.2.0/26", "2.3.4.0/24", "1.6.0.0/16", "1.2.2.64/26"},
		ASN:      1203874,
	}
	networks26 = append(networks26, network26)

	// no match -> no label
	nomatchNetworks := []*metal.MachineNetwork{}
	nomatchNetwork := &metal.MachineNetwork{
		Primary:  true,
		IPs:      []string{"10.2.0.4"},
		Prefixes: []string{"1.2.0.0/22", "1.2.2.0/22", "2.3.4.0/24"},
		ASN:      1203874,
	}
	nomatchNetworks = append(nomatchNetworks, nomatchNetwork)

	// no ip -> no label
	noipNetworks := []*metal.MachineNetwork{}
	noipNetwork := &metal.MachineNetwork{
		Primary:  true,
		IPs:      []string{},
		Prefixes: []string{"1.2.0.0/22", "1.2.2.0/22", "2.3.4.0/24"},
		ASN:      1203874,
	}
	noipNetworks = append(noipNetworks, noipNetwork)

	tests := []struct {
		name    string
		machine *metal.Machine
		want    []string
	}{
		{
			name: "simple",
			machine: &metal.Machine{
				Allocation: &metal.MachineAllocation{
					MachineNetworks: networks,
				},
				RackID: "rack01",
				IPMI: metal.IPMI{
					Fru: metal.Fru{
						ChassisPartSerial: "chassis123",
					},
				},
			},
			want: []string{
				"machine.metal-pod.io/network.primary.asn=1203874",
				"machine.metal-pod.io/network.primary.localbgp.ip=1.2.0.0",
				"machine.metal-pod.io/rack=rack01",
				"machine.metal-pod.io/chassis=chassis123",
			},
		},
		{
			name: "simple26",
			machine: &metal.Machine{
				Allocation: &metal.MachineAllocation{
					MachineNetworks: networks26,
				},
				RackID: "rack01",
				IPMI: metal.IPMI{
					Fru: metal.Fru{
						ChassisPartSerial: "chassis123",
					},
				},
			},
			want: []string{
				"machine.metal-pod.io/network.primary.asn=1203874",
				"machine.metal-pod.io/network.primary.localbgp.ip=1.2.2.64",
				"machine.metal-pod.io/rack=rack01",
				"machine.metal-pod.io/chassis=chassis123",
			},
		},
		{
			name: "ip does not match prefix",
			machine: &metal.Machine{
				Allocation: &metal.MachineAllocation{
					MachineNetworks: nomatchNetworks,
				},
				RackID: "rack01",
				IPMI: metal.IPMI{
					Fru: metal.Fru{
						ChassisPartSerial: "chassis123",
					},
				},
			},
			want: []string{
				"machine.metal-pod.io/network.primary.asn=1203874",
				"machine.metal-pod.io/rack=rack01",
				"machine.metal-pod.io/chassis=chassis123",
			},
		},
		{
			name: "no ip",
			machine: &metal.Machine{
				Allocation: &metal.MachineAllocation{
					MachineNetworks: nomatchNetworks,
				},
				RackID: "rack01",
				IPMI: metal.IPMI{
					Fru: metal.Fru{
						ChassisPartSerial: "chassis123",
					},
				},
			},
			want: []string{
				"machine.metal-pod.io/network.primary.asn=1203874",
				"machine.metal-pod.io/rack=rack01",
				"machine.metal-pod.io/chassis=chassis123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := additionalTags(tt.machine); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("additionalTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
