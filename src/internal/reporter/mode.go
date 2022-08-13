package reporter

import (
	"strconv"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype/primefimp"
	log "github.com/sirupsen/logrus"
)

const (
	eepCgxPath          = "/cgx/eepl.cgx"
	iosPath             = "/html_old/cgx/ios.cgx?ios=0"
	backlightEepAddress = "EEV0002"
)

type mode struct {
	host   string
	client doer
}

// NewMode creates a new mode reporter.
func NewMode(host string, client doer) Rest {
	return &mode{host: host, client: client}
}

// Report waits for vinculum mode change message and sets the wall controllers light.
// For home mode it sets the light to 40.
// For other modes it sets the light to 0.
func (m *mode) Report(msg *fimpgo.Message) error {
	var cmd primefimp.Notify

	err := msg.Payload.GetObjectValue(&cmd)
	if err != nil {
		return err
	}

	if cmd.Cmd != "set" || cmd.Component != "hub" || cmd.Id != "mode" {
		return nil
	}

	modeChange := cmd.GetModeChange()
	if modeChange == nil {
		log.Warn("Mode change message is not valid")

		return err
	}

	if modeChange.Current == "home" {
		return m.publish(40)
	}

	return m.publish(0)
}

func (m *mode) publish(v int) error {
	err := publish(m.client, m.host+eepCgxPath, backlightEepAddress, strconv.Itoa(v))
	if err != nil {
		return err
	}

	// publishing two changes as we don't know what's the current status. It's faster than checking it before sending.
	err = publish(m.client, m.host+iosPath, "IOO0122", "0")
	if err != nil {
		return err
	}

	return publish(m.client, m.host+iosPath, "IOO0122", "10")
}
