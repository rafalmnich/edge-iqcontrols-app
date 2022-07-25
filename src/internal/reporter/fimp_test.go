package reporter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/futurehomeno/fimpgo"
	"github.com/stretchr/testify/require"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/transformer"
)

func TestFimp_Report(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		add     int64
		val     int64
		want    *fimpgo.Message
		wantErr bool
	}{
		{
			name: "Report for luminosity",
			add:  1,
			val:  100,
			want: &fimpgo.Message{
				Addr:    sampleAddress(t, 1),
				Payload: fimpgo.NewIntMessage("evt.sensor.report", "sensor_lumin", 100, nil, nil, nil),
			},
		},
		{
			name:    "device not found",
			add:     404,
			val:     100,
			wantErr: true,
		},
	}

	mqtt := GetMQTT(t)

	t.Cleanup(mqtt.Stop)
	msgChan := Subscribe(t, mqtt, "pt:j1/mt:evt/rt:dev/rn:iqcontrols/ad:1/#")

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			checked := make(chan struct{})

			if !tt.wantErr {
				go check(t, msgChan, checked, checkSensorLuminReport)
			}

			f := reporter.NewFimp(mqtt, transformer.NewFimp([]config.Device{
				{
					Name:        "light in office",
					Address:     2,
					ServiceName: "out_bin_switch",
				},
				{
					Name:        "north luminance sensor",
					Address:     1,
					ServiceName: "sensor_lumin",
				},
			}))

			err := f.Report(tt.add, tt.val)
			if tt.wantErr {
				require.Error(t, err)
				close(checked)
			} else {
				require.NoError(t, err)
			}

			<-checked
		})
	}
}

// Subscribe subscribes to the given topic and returns a channel with messages.
func Subscribe(t *testing.T, mqtt *fimpgo.MqttTransport, topic string) chan *fimpgo.Message {
	t.Helper()

	err := mqtt.Subscribe(topic)
	require.NoError(t, err)

	msgChan := make(chan *fimpgo.Message, 10)
	mqtt.RegisterChannel("response_channel", msgChan)

	return msgChan
}

// GetMQTT returns a new MQTT transport.
func GetMQTT(t *testing.T) *fimpgo.MqttTransport {
	mqtt := fimpgo.NewMqttTransport(
		"tcp://localhost:11883",
		"app_tests",
		"",
		"",
		true,
		1,
		1,
	)

	err := mqtt.Start()
	require.NoError(t, err)

	return mqtt
}

func check(t *testing.T, msgChan chan *fimpgo.Message, checked chan struct{}, checkFunc func(t *testing.T, msg *fimpgo.Message)) {
	t.Helper()

	defer close(checked)

	select {
	case msg := <-msgChan:
		checkFunc(t, msg)

		return
	case <-time.After(time.Hour):
		t.Fatal("timeout")
	}

}

func sampleAddress(t *testing.T, add int64) *fimpgo.Address {
	t.Helper()

	a, err := fimpgo.NewAddressFromString(fmt.Sprintf("pt:j1/mt:evt/rt:dev/rn:iqcontrols/ad:1/sv:sensor_lumin/ad:%d_0", add))
	require.NoError(t, err)

	return a
}

func checkSensorLuminReport(t *testing.T, msg *fimpgo.Message) {
	require.Equal(t, sampleAddress(t, 1), msg.Addr)
	require.Equal(t, "evt.sensor.report", msg.Payload.Type)
	require.Equal(t, "int", msg.Payload.ValueType)
	require.Equal(t, "sensor_lumin", msg.Payload.Service)
	require.Equal(t, "Lux", msg.Payload.Properties["unit"])

	value, err := msg.Payload.GetIntValue()
	require.NoError(t, err)
	require.Equal(t, int64(100), value)
}
