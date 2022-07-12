package transformer

import (
	"fmt"

	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

// Mapper represents a value and topic fimp for specific devices.
type Mapper interface {
	Message(add int64, val int64) (*fimpgo.Message, error)
}

// Strategy represents a transformer strategy with fimp and device strategies.
type Strategy interface {
	FimpStrategy
	DeviceStrategy
}

// FimpStrategy represents a strategy for a specific device type.
type FimpStrategy interface {
	FimpMessage(device config.Device, val int64) *fimpgo.FimpMessage
}

// DeviceStrategy represents a strategy for mapping device topic and value to right address and value for mass driver.
type DeviceStrategy interface {
	Value(device config.Device, msg *fimpgo.Message) (value string, err error)
}

// FimpStrategyFunc represents a function that maps address and value to fimp message.
type FimpStrategyFunc func(device config.Device, val int64) *fimpgo.FimpMessage

// FimpMessage implements Mapper.FimpMessage.
func (m FimpStrategyFunc) FimpMessage(d config.Device, v int64) *fimpgo.FimpMessage {
	return m(d, v)
}

type fimp struct {
	devices    []config.Device
	strategies map[string]FimpStrategy
}

// NewFimp creates a new value and topic fimp with different strategies for each device type.
func NewFimp(devices []config.Device) Mapper {
	return &fimp{
		devices: devices,
		strategies: map[string]FimpStrategy{
			SensorLumin:    FimpStrategyFunc(LuminanceFimpMessage),
			SensorPresence: FimpStrategyFunc(PresenceFimpMessage),
			SensorTemp:     FimpStrategyFunc(TemperatureFimpMessage),
			OutBinSwitch:   FimpStrategyFunc(BinSwitchFimpMessage),
			OutLvlSwitch:   FimpStrategyFunc(LvlSwitchFimpMessage),
		},
	}
}

// Message sends a value to mqtt.
func (f *fimp) Message(add int64, val int64) (*fimpgo.Message, error) {
	d := f.findDevice(add)
	if d.Address == 0 {
		return nil, fmt.Errorf("Device with address %d not found", add)
	}

	strategy, ok := f.strategies[d.ServiceName]
	if !ok {
		return nil, fmt.Errorf("strategy for Device type %s not found", d.ServiceName)
	}

	return &fimpgo.Message{
		Addr:    f.address(d),
		Payload: strategy.FimpMessage(d, val),
	}, nil
}

func (f *fimp) address(d config.Device) *fimpgo.Address {
	return &fimpgo.Address{
		PayloadType:     "j1",
		MsgType:         "evt",
		ResourceType:    "dev",
		ResourceName:    "iqcontrols",
		ResourceAddress: "1",
		ServiceName:     d.ServiceName,
		ServiceAddress:  fmt.Sprintf("%d_0", d.Address),
	}
}

func (f *fimp) findDevice(add int64) config.Device {
	for _, d := range f.devices {
		if d.Address == add {
			return d
		}
	}

	return config.Device{}
}
