package transformer

import (
	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

// SensorPresence represents a strategy for presence sensor.
const SensorPresence = "sensor_presence"

// PresenceFimpMessage returns a fimp event for presence sensor.
func PresenceFimpMessage(d config.Device, val int64) *fimpgo.FimpMessage {
	min := boolMin(d.Config["minValue"])

	return fimpgo.NewBoolMessage(
		sensorReport,
		SensorPresence,
		val > min,
		nil,
		nil,
		nil,
	)
}

func boolMin(v interface{}) int64 {
	min, ok := v.(int)
	if !ok {
		return 0
	}

	return int64(min)
}
