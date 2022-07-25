package transformer

import (
	"math"
	"strconv"

	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

const (
	// OutLvlSwitch represents a strategy for out lvl switch.
	OutLvlSwitch = "out_lvl_switch"

	levelReport = "evt.lvl.report"
)

// LvlSwitchFimpMessage returns a fimp event for lvl switch.
func LvlSwitchFimpMessage(device config.Device, val int64) *fimpgo.FimpMessage {
	return fimpgo.NewIntMessage(
		levelReport,
		OutLvlSwitch,
		applyIntMultiplier(device.Config, val),
		nil,
		nil,
		nil,
	)
}

// LvlSwitchValue returns value for lvl switch.
func LvlSwitchValue(device config.Device, msg *fimpgo.Message) (value string, err error) {
	v, err := val(msg)
	if err != nil {
		return "", err
	}

	return applyReverseMultiplier(device.Config, float64(v)), nil
}

func val(msg *fimpgo.Message) (int64, error) {
	if msg.Payload.Type != "cmd.binary.set" {
		return msg.Payload.GetIntValue()
	}

	v, err := msg.Payload.GetBoolValue()
	if err != nil {
		return 0, err
	}

	if v {
		return 100, nil
	}

	return 0, nil
}

func applyReverseMultiplier(c map[string]interface{}, val float64) string {
	multiplier := multiplier(c["multiplier"])

	return strconv.Itoa(int(val / multiplier))
}

func multiplier(m interface{}) float64 {
	multiplier, ok := m.(float64)
	if !ok {
		return 1
	}

	return multiplier
}

func applyIntMultiplier(config map[string]interface{}, val int64) int64 {
	multiplier := multiplier(config["multiplier"])

	v := float64(val) * multiplier
	v = math.Round(v*10) / 10

	return int64(v)
}
