package transformer

import (
	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

const (
	// OutBinSwitch represents a strategy for out bin switch.
	OutBinSwitch = "out_bin_switch"

	binaryReport = "evt.binary.report"
)

// BinSwitchFimpMessage returns a fimp event for bin switch.
func BinSwitchFimpMessage(_ config.Device, val int64) *fimpgo.FimpMessage {
	return fimpgo.NewBoolMessage(
		binaryReport,
		OutBinSwitch,
		val > 0,
		nil,
		nil,
		nil,
	)
}

// BinSwitchValue returns value for bin switch.
// for true - returns "1000"
// for false - returns "0"
func BinSwitchValue(_ config.Device, msg *fimpgo.Message) (value string, err error) {
	val, err := msg.Payload.GetBoolValue()
	if err != nil {
		return "", err
	}

	if val {
		return "1000", nil
	}

	return "0", nil
}
