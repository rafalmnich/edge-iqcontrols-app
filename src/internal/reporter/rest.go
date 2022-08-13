package reporter

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

type (
	// RestPublisher represents a rest publisher.
	RestPublisher interface {
		Publish(device config.Device, value string) error
	}
	// DeviceMapper represents a device mapper.
	DeviceMapper interface {
		Device(msg *fimpgo.Message) (device config.Device, value string, err error)
	}
	// Rest represents a rest reporter.
	Rest interface {
		Report(msg *fimpgo.Message) error
	}

	doer interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

type rest struct {
	publisher RestPublisher
	mapper    DeviceMapper
}

// NewRest creates a new rest reporter.
func NewRest(publisher RestPublisher, mapper DeviceMapper) Rest {
	return &rest{publisher: publisher, mapper: mapper}
}

// Report reports a value taken from fimp message to mass rest api.
func (r *rest) Report(msg *fimpgo.Message) error {
	device, value, err := r.mapper.Device(msg)
	if err != nil {
		return err
	}

	return r.publisher.Publish(device, value)
}

type restPublisher struct {
	client         doer
	massLightsCgx  string
	massHeatingCgx string
}

// NewRestPublisher creates a new rest publisher.
func NewRestPublisher(client doer, massLightCgx string, massHeatingCgx string) RestPublisher {
	return &restPublisher{client: client, massLightsCgx: massLightCgx, massHeatingCgx: massHeatingCgx}
}

// Publish sends a value to mass rest api.
func (r *restPublisher) Publish(device config.Device, value string) error {
	cgxPath := r.massLightsCgx

	if device.Mass == "heating" {
		cgxPath = r.massHeatingCgx
	}

	return publish(r.client, cgxPath, device.VariableName, value)
}

func publish(client doer, cgxPath, address, value string) error {
	form := url.Values{}
	form.Add(address, value)

	req, err := http.NewRequest(http.MethodPost, cgxPath, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Debugf("Sending request: %s, data: %s", req.URL.String(), form.Encode())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
