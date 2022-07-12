package tests

import (
	"fmt"
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/stretchr/testify/assert"
)

// AssertMessage checks if all important fields are equal in Message.
func AssertMessage(t *testing.T, expected *fimpgo.Message, got *fimpgo.Message) {
	t.Helper()

	assert.Equal(t, expected.Addr, got.Addr)
	assert.Equal(t, expected.Payload.Type, got.Payload.Type)
	assert.Equal(t, expected.Payload.Service, got.Payload.Service)
	assert.Equal(t, expected.Payload.ValueType, got.Payload.ValueType)
	assert.Equal(t, expected.Payload.Value, got.Payload.Value)
	assert.Equal(t, expected.Payload.Properties, got.Payload.Properties)
}

// DeviceAddress returns device address for given device address.
func DeviceAddress(addr int) *fimpgo.Address {
	return &fimpgo.Address{
		PayloadType:     "j1",
		MsgType:         "cmd",
		ResourceType:    "dev",
		ResourceName:    "iqcontrols",
		ResourceAddress: "1",
		ServiceName:     "out_bin_switch",
		ServiceAddress:  serviceAddr(addr),
	}
}

func serviceAddr(addr int) string {
	return fmt.Sprintf("%d_0", addr)
}
