package reporter

import (
	"github.com/futurehomeno/fimpgo"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/transformer"
)

// Reporter represents a values reporter for specific devices.
type Reporter interface {
	Report(add int64, val int64) error
}

// FimpPublisher represents a fimpgo publisher.
type FimpPublisher interface {
	Publish(addr *fimpgo.Address, fimpMsg *fimpgo.FimpMessage) error
}

type fimp struct {
	publisher FimpPublisher
	mapper    transformer.Mapper
}

// NewFimp creates a new fimp reporter.
func NewFimp(mqtt FimpPublisher, mapper transformer.Mapper) Reporter {
	return &fimp{publisher: mqtt, mapper: mapper}
}

// Report reports a value to mqtt
func (f *fimp) Report(add int64, val int64) error {
	msg, err := f.mapper.Message(add, val)
	if err != nil {
		return err
	}

	return f.publisher.Publish(msg.Addr, msg.Payload)
}
