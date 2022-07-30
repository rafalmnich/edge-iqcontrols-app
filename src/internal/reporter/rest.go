package reporter

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
)

const cgxPath = "/cgx/custom/all.cgx"

type (
	// RestPublisher represents a rest publisher.
	RestPublisher interface {
		Publish(address, value string) error
	}
	// DeviceMapper represents a device mapper.
	DeviceMapper interface {
		Device(msg *fimpgo.Message) (address string, value string, err error)
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
	address, value, err := r.mapper.Device(msg)
	if err != nil {
		return err
	}

	// todo: publish also a fimp event after this. ??
	return r.publisher.Publish(address, value)
}

type restPublisher struct {
	host   string
	client doer
}

// NewRestPublisher creates a new rest publisher.
func NewRestPublisher(host string, client doer) RestPublisher {
	return &restPublisher{host: host, client: client}
}

// Publish sends a value to mass rest api.
func (r *restPublisher) Publish(address, value string) error {
	return publish(r.client, r.host, cgxPath, address, value)
}

func publish(client doer, host, cgxPath, address, value string) error {
	form := url.Values{}
	form.Add(address, value)

	req, err := http.NewRequest(http.MethodPost, host+cgxPath, strings.NewReader(form.Encode()))
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
