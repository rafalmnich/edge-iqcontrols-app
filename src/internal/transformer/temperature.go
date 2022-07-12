package transformer

import (
	"math"

	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

// SensorTemp represents a strategy for temperature sensor.
const SensorTemp = "sensor_temp"

// TemperatureFimpMessage returns a fimp event for temperature sensor.
func TemperatureFimpMessage(d config.Device, val int64) *fimpgo.FimpMessage {
	return fimpgo.NewFloatMessage(
		sensorReport,
		SensorTemp,
		applyFloatMultiplier(d.Config, val),
		fimpgo.Props{
			"unit": "C",
		},
		nil,
		nil,
	)
}

func applyFloatMultiplier(config map[string]interface{}, val int64) float64 {
	multiplier := multiplier(config["multiplier"])

	v := float64(val) * multiplier

	return math.Round(v*10) / 10
}
