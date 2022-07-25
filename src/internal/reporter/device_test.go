package reporter_test

import (
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	. "github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
)

func TestInclusionReport(t *testing.T) {
	tests := []struct {
		name    string
		device  config.Device
		want    fimptype.ThingInclusionReport
		wantErr bool
	}{
		{
			name: "sensor lumin",
			device: config.Device{
				Name:        "North luminance",
				Address:     6,
				ServiceName: "sensor_lumin",
				Type:        "in",
				MsgType:     "evt.sensor.report",
				ValueType:   "float",
			},
			want: sensorLuminReport(),
		},
		{
			name: "sensor presence",
			device: config.Device{
				Name:        "Bathroom small presence sensor",
				Address:     166,
				ServiceName: "sensor_presence",
				Type:        "in",
				MsgType:     "evt.sensor.report",
				ValueType:   "bool",
			},
			want: sensorPresenceReport(),
		},
	}

	mqtt := GetMQTT(t)
	t.Cleanup(mqtt.Stop)

	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			msgChan := Subscribe(t, mqtt, "pt:j1/mt:evt/rt:ad/rn:flow/ad:1_0")
			t.Cleanup(func() {
				close(msgChan)
			})

			checked := make(chan struct{})

			go check(t, msgChan, checked, checkInclusionReport(tt.want))

			inclusion := NewDevice(mqtt)

			err := inclusion.InclusionReport(tt.device)

			assert.NoError(t, err)

			<-checked
		})
	}
}

func sensorLuminReport() fimptype.ThingInclusionReport {
	return fimptype.ThingInclusionReport{
		IntegrationId:  "iqcontrols",
		Address:        "6",
		Type:           "",
		ProductHash:    "iq_hash_6",
		CommTechnology: "wire",
		ProductId:      "iq_6",
		ProductName:    "North luminance",
		ManufacturerId: "iqcontrols",
		DeviceId:       "6",
		PowerSource:    "ac",
		Groups:         []string{"1"},
		Services: []fimptype.Service{
			{
				Name:    "sensor_lumin",
				Alias:   "North luminance",
				Address: "/rt:dev/rn:iqcontrols/ad:1/sv:sensor_lumin/ad:6_0",
				Enabled: true,
				Groups:  []string{"1"},
				Interfaces: []fimptype.Interface{
					{
						Type:      "in",
						MsgType:   "evt.sensor.report",
						ValueType: "float",
						Version:   "1",
					},
				},
			},
		},
	}
}

func sensorPresenceReport() fimptype.ThingInclusionReport {
	return fimptype.ThingInclusionReport{
		IntegrationId:  "iqcontrols",
		Address:        "166",
		Type:           "",
		ProductHash:    "iq_hash_166",
		CommTechnology: "wire",
		ProductId:      "iq_166",
		ProductName:    "Bathroom small presence sensor",
		ManufacturerId: "iqcontrols",
		DeviceId:       "166",
		PowerSource:    "ac",
		Groups:         []string{"1"},
		Services: []fimptype.Service{
			{
				Name:    "sensor_presence",
				Alias:   "Bathroom small presence sensor",
				Address: "/rt:dev/rn:iqcontrols/ad:1/sv:sensor_presence/ad:166_0",
				Enabled: true,
				Groups:  []string{"1"},
				Interfaces: []fimptype.Interface{
					{
						Type:      "in",
						MsgType:   "evt.sensor.report",
						ValueType: "bool",
						Version:   "1",
					},
				},
			},
		},
	}
}

func checkInclusionReport(r fimptype.ThingInclusionReport) func(t *testing.T, msg *fimpgo.Message) {
	return func(t *testing.T, msg *fimpgo.Message) {
		var report fimptype.ThingInclusionReport

		err := msg.Payload.GetObjectValue(&report)
		require.NoError(t, err)

		assert.Equal(t, r, report)
		assert.Equal(t, "evt.thing.inclusion_report", msg.Payload.Type)
		assert.Equal(t, "iqcontrols", msg.Payload.Service)
	}
}
