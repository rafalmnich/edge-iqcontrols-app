package transformer

import (
	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

const (
	// SensorLumin represents a strategy for luminance sensor.
	SensorLumin = "sensor_lumin"

	sensorReport = "evt.sensor.report"
)

// LuminanceFimpMessage returns a fimp event for luminance sensor.
func LuminanceFimpMessage(_ config.Device, val int64) *fimpgo.FimpMessage {
	return fimpgo.NewFloatMessage(
		sensorReport,
		SensorLumin,
		float64(val),
		nil,
		nil,
		nil,
	)
}
