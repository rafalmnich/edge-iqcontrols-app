package transformer_test

import (
	"fmt"
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/stretchr/testify/require"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/tests"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/transformer"
)

func TestFimp_Message(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		devices []config.Device
		val     int64
		add     int64
		want    *fimpgo.Message
		wantErr bool
	}{
		{
			name:    "report for temperature",
			devices: sampleDevices(t),
			val:     123,
			add:     4,
			want: &fimpgo.Message{
				Addr: sensorTempAddr(t, 4),
				Payload: fimpgo.NewFloatMessage(
					"evt.sensor.report",
					"sensor_temp",
					12.3,
					fimpgo.Props{"unit": "C"},
					nil,
					nil,
				),
			},
		},
		{
			name:    "report for temperature, but device not found",
			devices: sampleDevices(t),
			val:     123,
			add:     404,
			wantErr: true,
		},
		{
			name:    "report for temperature, device type not found",
			devices: sampleDevices(t),
			val:     123,
			add:     666,
			wantErr: true,
		},
	}

	for _, ttt := range testCases {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mapper := transformer.NewFimp(tt.devices)
			got, err := mapper.Message(tt.add, tt.val)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tests.AssertMessage(t, tt.want, got)
		})
	}
}

func sampleDevices(t *testing.T) []config.Device {
	t.Helper()

	return []config.Device{
		{
			Name:        "luminosity",
			Address:     3,
			ServiceName: "sensor_lumin",
			Config: map[string]interface{}{
				"minValue": 10,
			},
		},
		{
			Name:        "temp sensor",
			Address:     4,
			ServiceName: "sensor_temp",
			Config: map[string]interface{}{
				"multiplier": 0.1,
			},
		},
		{
			Name:         "bin switch",
			Address:      5,
			ServiceName:  "out_bin_switch",
			VariableName: "IOO005",
		},
		{
			Name:         "level switch",
			Address:      6,
			ServiceName:  "out_lvl_switch",
			VariableName: "IOO006",
			Config: map[string]interface{}{
				"multiplier": 0.1,
			},
		},
		{
			Name:         "level switch - no multiplier",
			Address:      7,
			ServiceName:  "out_lvl_switch",
			VariableName: "IOO007",
		},
		{
			Name:        "unknown device",
			Address:     666,
			ServiceName: "unknown_type",
		},
	}
}

func sensorTempAddr(t *testing.T, addr int) *fimpgo.Address {
	t.Helper()

	address, err := fimpgo.NewAddressFromString(fmt.Sprintf("pt:j1/mt:evt/rt:dev/rn:iqcontrols/ad:1/sv:sensor_temp/ad:%d_0", addr))
	require.NoError(t, err)

	return address
}
