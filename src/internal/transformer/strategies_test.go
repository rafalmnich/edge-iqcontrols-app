package transformer_test

import (
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/transformer"
)

func TestFimpStrategies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		strategy transformer.FimpStrategy
		val      int64
		device   config.Device
		want     *fimpgo.FimpMessage
	}{
		{
			name:     "report for luminosity",
			strategy: transformer.FimpStrategyFunc(transformer.LuminanceFimpMessage),
			val:      100,
			want: &fimpgo.FimpMessage{
				Type:      "evt.sensor.report",
				Service:   "sensor_lumin",
				ValueType: "float",
				Value:     float64(100),
			},
		},
		{
			name:     "report for presence",
			strategy: transformer.FimpStrategyFunc(transformer.PresenceFimpMessage),
			val:      12,
			device:   config.Device{},
			want:     presence(t, true),
		},
		{
			name:     "report for presence with minValue - on",
			strategy: transformer.FimpStrategyFunc(transformer.PresenceFimpMessage),
			val:      11,
			device: config.Device{
				Config: map[string]interface{}{
					"minValue": 10,
				},
			},
			want: presence(t, true),
		},
		{
			name:     "report for presence with minValue - off",
			strategy: transformer.FimpStrategyFunc(transformer.PresenceFimpMessage),
			val:      10,
			device: config.Device{
				Config: map[string]interface{}{
					"minValue": 10,
				},
			},
			want: presence(t, false),
		},
		{
			name:     "report for temperature",
			strategy: transformer.FimpStrategyFunc(transformer.TemperatureFimpMessage),
			val:      234,
			device: config.Device{
				Config: map[string]interface{}{
					"multiplier": 0.1,
				},
			},
			want: &fimpgo.FimpMessage{
				Type:      "evt.sensor.report",
				Service:   "sensor_temp",
				ValueType: "float",
				Value:     23.4,
				Properties: map[string]string{
					"unit": "C",
				},
			},
		},
		{
			name:     "report for binary switch on",
			strategy: transformer.FimpStrategyFunc(transformer.BinSwitchFimpMessage),
			val:      1000,
			device:   config.Device{},
			want: &fimpgo.FimpMessage{
				Type:      "evt.binary.report",
				Service:   "out_bin_switch",
				ValueType: "bool",
				Value:     true,
			},
		},
		{
			name:     "report for binary switch off",
			strategy: transformer.FimpStrategyFunc(transformer.BinSwitchFimpMessage),
			val:      0,
			device:   config.Device{},
			want: &fimpgo.FimpMessage{
				Type:      "evt.binary.report",
				Service:   "out_bin_switch",
				ValueType: "bool",
				Value:     false,
			},
		},
		{
			name:     "report for level switch 37%",
			strategy: transformer.FimpStrategyFunc(transformer.LvlSwitchFimpMessage),
			val:      370,
			device: config.Device{
				Config: map[string]interface{}{
					"multiplier": 0.1,
				},
			},
			want: &fimpgo.FimpMessage{
				Type:      "evt.level.report",
				Service:   "out_lvl_switch",
				ValueType: "int",
				Value:     int64(37),
			},
		},
	}

	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.NotNil(t, tt.strategy, "define strategy for this test case!")

			got := tt.strategy.FimpMessage(tt.device, tt.val)

			assert.Equal(t, tt.want.Service, got.Service)
			assert.Equal(t, tt.want.Value, got.Value)
			assert.Equal(t, tt.want.ValueType, got.ValueType)
			assert.Equal(t, tt.want.Type, got.Type)
		})
	}
}

func presence(t *testing.T, v bool) *fimpgo.FimpMessage {
	t.Helper()

	return &fimpgo.FimpMessage{
		Type:      "evt.sensor.report",
		Service:   "sensor_presence",
		ValueType: "bool",
		Value:     v,
	}
}
