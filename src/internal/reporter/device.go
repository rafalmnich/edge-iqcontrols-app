package reporter

import (
	"encoding/json"
	"strconv"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

// Device represents a device device reporter.
type Device interface {
	InclusionReport(device config.Device) error
	ExclusionReport(deviceID int) error
}

type id int

type device struct {
	publisher FimpPublisher
	devices   map[id]config.Device
}

// NewDevice creates a new device reporter.
func NewDevice(publisher FimpPublisher) Device {
	return &device{publisher: publisher}
}

// InclusionReport reports an inclusion of a device.
func (i *device) InclusionReport(device config.Device) error {
	inclusion := i.inclusion(device)

	msg := fimpgo.NewObjectMessage("evt.thing.inclusion_report", "iqcontrols", inclusion, nil, nil, nil)
	addr, _ := fimpgo.NewAddressFromString("pt:j1/mt:evt/rt:ad/rn:flow/ad:1")

	b, _ := json.Marshal(msg)
	_ = b

	return i.publisher.Publish(addr, msg)
}

// ExclusionReport removes device from
func (i *device) ExclusionReport(_ int) error {
	// TODO implement me
	panic("implement me")
}

func (i *device) inclusion(d config.Device) *fimptype.ThingInclusionReport {
	return &fimptype.ThingInclusionReport{
		IntegrationId:  "iqcontrols",
		Address:        address(d.Address),
		Type:           "",
		ProductHash:    "iq_hash_" + address(d.Address),
		CommTechnology: "wire",
		ProductId:      "iq_" + address(d.Address),
		ProductName:    d.Name,
		ManufacturerId: "iqcontrols",
		DeviceId:       address(d.Address),
		PowerSource:    "ac",
		Groups:         []string{"1"},
		Services: []fimptype.Service{
			{
				Name:    d.ServiceName,
				Alias:   d.Name,
				Address: "/rt:dev/rn:iqcontrols/ad:1/sv:" + d.ServiceName + "/ad:" + addressZero(d.Address),
				Enabled: true,
				Groups:  []string{"1"},
				Interfaces: []fimptype.Interface{
					{
						Type:      "out",
						MsgType:   d.MsgType,
						ValueType: d.ValueType,
						Version:   "1",
					},
				},
			},
		},
	}
}

func addressZero(a int64) string {
	return address(a) + "_0"
}

func address(a int64) string {
	return strconv.Itoa(int(a))
}
