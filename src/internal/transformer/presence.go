package transformer

import (
	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

const (
	// SensorPresence represents a strategy for presence sensor.
	SensorPresence       = "sensor_presence"
	sensorPresenceReport = "evt.presence.report"
)

// PresenceFimpMessage returns a fimp event for presence sensor.
func PresenceFimpMessage(d config.Device, val int64) *fimpgo.FimpMessage {
	min := boolMin(d.Config["minValue"])

	return fimpgo.NewBoolMessage(
		sensorPresenceReport,
		SensorPresence,
		val > min,
		nil,
		nil,
		nil,
	)
}

func boolMin(v interface{}) int64 {
	min, ok := v.(float64)
	if !ok {
		return 0
	}

	return int64(min)
}
