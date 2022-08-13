package transformer_test

import (
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/stretchr/testify/assert"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/tests"
	. "github.com/rafalmnich/edge-iqcontrols-app/internal/transformer"
)

func TestDevice_Device(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		msg     *fimpgo.Message
		addr    string
		value   string
		wantErr bool
	}{
		{
			name:  "check bin switch addr 5, value 1000",
			msg:   binSwitchCommand(5, true),
			addr:  "IOO005",
			value: "1000",
		},
		{
			name:  "check bin switch addr 5, value 0",
			msg:   binSwitchCommand(5, false),
			addr:  "IOO005",
			value: "0",
		},
		{
			name:  "level switch addr 6, value 75%",
			msg:   levelSwitchCommand(6, 75),
			addr:  "IOO006",
			value: "750",
		},
		{
			name:    "device not found by that address",
			msg:     binSwitchCommand(404, true),
			wantErr: true,
		},
		{
			name:    "device type found but of unexpected type",
			msg:     unknownTypeCommand(),
			wantErr: true,
		},
		{
			name:    "bad type command",
			msg:     wrongBinSwitchCmd(5),
			wantErr: true,
		},
		{
			name:    "level switch bad command",
			msg:     wrongLvlSwitchCmd(6),
			wantErr: true,
		},
		{
			name:  "level switch turn off",
			msg:   lvlSwitchTurnOff(6),
			addr:  "IOO006",
			value: "0",
		},
		{
			name:  "level switch turn on",
			msg:   lvlSwitchTurnOn(6),
			addr:  "IOO006",
			value: "1000",
		},
		{
			name:  "level switch no multiplier in config",
			msg:   levelSwitchCommand(7, 75),
			addr:  "IOO007",
			value: "75",
		},
	}

	for _, ttt := range testCases {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDevice(sampleDevices(t))

			device, value, err := d.Device(tt.msg)
			if tt.wantErr {
				assert.Error(t, err)

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.addr, device.VariableName)
			assert.Equal(t, tt.value, value)
		})
	}
}

func unknownTypeCommand() *fimpgo.Message {
	msg := fimpgo.NewBoolMessage("cmd.binary.set", "unknown_type", true, nil, nil, nil)

	a := &fimpgo.Address{
		PayloadType:     "j1",
		MsgType:         "cmd",
		ResourceType:    "dev",
		ResourceName:    "iqcontrols",
		ResourceAddress: "1",
		ServiceName:     "out_bin_switch",
		ServiceAddress:  "666_0",
	}

	return &fimpgo.Message{
		Addr:    a,
		Payload: msg,
	}
}

func binSwitchCommand(addr int, val bool) *fimpgo.Message {
	msg := fimpgo.NewBoolMessage("cmd.binary.set", "out_bin_switch", val, nil, nil, nil)

	return &fimpgo.Message{
		Addr:    tests.DeviceAddress(addr),
		Payload: msg,
	}
}

func levelSwitchCommand(addr int, val int64) *fimpgo.Message {
	msg := fimpgo.NewIntMessage("cmd.lvl.set", "out_lvl_switch", val, nil, nil, nil)

	return &fimpgo.Message{
		Addr:    tests.DeviceAddress(addr),
		Payload: msg,
	}
}

func wrongLvlSwitchCmd(addr int) *fimpgo.Message {
	msg := fimpgo.NewBoolMessage("cmd.lvl.set", "out_lvl_switch", true, nil, nil, nil)

	return &fimpgo.Message{
		Addr:    tests.DeviceAddress(addr),
		Payload: msg,
	}
}

func lvlSwitchTurnOff(addr int) *fimpgo.Message {
	msg := fimpgo.NewBoolMessage("cmd.binary.set", "out_lvl_switch", false, nil, nil, nil)

	return &fimpgo.Message{
		Addr:    tests.DeviceAddress(addr),
		Payload: msg,
	}
}

func lvlSwitchTurnOn(addr int) *fimpgo.Message {
	msg := lvlSwitchTurnOff(addr)

	msg.Payload.Value = true

	return msg
}

func wrongBinSwitchCmd(addr int) *fimpgo.Message {
	msg := fimpgo.NewIntMessage("cmd.binary.set", "out_bin_switch", 1, nil, nil, nil)

	return &fimpgo.Message{
		Addr:    tests.DeviceAddress(addr),
		Payload: msg,
	}

}
