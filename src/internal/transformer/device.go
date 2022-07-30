package transformer

import (
	"fmt"
	"strconv"

	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

// DeviceStrategyFunc is a function type used for strategy.
type DeviceStrategyFunc func(device config.Device, msg *fimpgo.Message) (value string, err error)

// Value returns a value mapped from a given fimp message.
func (d DeviceStrategyFunc) Value(device config.Device, msg *fimpgo.Message) (value string, err error) {
	return d(device, msg)
}

// Device represents a value and topic fimp for specific devices.
type Device struct {
	devices    []config.Device
	strategies map[string]DeviceStrategy
}

// NewDevice creates a new device fimp with different strategies for each device type.
func NewDevice(devices []config.Device) *Device {
	return &Device{
		devices: devices,
		strategies: map[string]DeviceStrategy{
			OutBinSwitch: DeviceStrategyFunc(BinSwitchValue),
			OutLvlSwitch: DeviceStrategyFunc(LvlSwitchValue),
		},
	}
}

// Device returns a device address and name for a given fimp message.
func (r *Device) Device(msg *fimpgo.Message) (address string, value string, err error) {
	d := r.findDevice(msg)

	if d.Address == 0 {
		return "", "", fmt.Errorf("device with address %d not found", d.Address)
	}

	strategy, ok := r.strategies[d.ServiceName]
	if !ok {
		return "", "", fmt.Errorf("strategy for Device type %s not found", d.ServiceName)
	}

	value, err = strategy.Value(d, msg)
	if err != nil {
		return "", "", err
	}

	return d.VariableName, value, nil
}

func (r *Device) findDevice(msg *fimpgo.Message) config.Device {
	for _, d := range r.devices {
		address := strconv.Itoa(int(d.Address)) + "_0"

		if address == msg.Addr.ServiceAddress && d.ServiceName == msg.Payload.Service {
			return d
		}
	}

	return config.Device{}
}
