package reporter_test

import (
	"fmt"
	"testing"

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
			name: "report for luminosity",
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

	mqtt := fimpgo.NewMqttTransport(
		"tcp://localhost:11883",
		"app_tests",
		"guest",
		"guest",
		true,
		1,
		1,
	)
	err := mqtt.Start()
	require.NoError(t, err)

	t.Cleanup(mqtt.Stop)

	err = mqtt.Subscribe("pt:j1/mt:evt/rt:dev/rn:iqcontrols/ad:1/#")
	require.NoError(t, err)

	msgChan := make(chan *fimpgo.Message, 10)
	mqtt.RegisterChannel("response_channel", msgChan)

	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			checked := make(chan struct{})

			if !tt.wantErr {
				go check(t, msgChan, checked)
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

func check(t *testing.T, msgChan chan *fimpgo.Message, checked chan struct{}) {
	t.Helper()

	defer close(checked)

	for {
		select {
		case msg := <-msgChan:
			require.Equal(t, sampleAddress(t, 1), msg.Addr)
			require.Equal(t, "evt.sensor.report", msg.Payload.Type)
			require.Equal(t, "sensor_lumin", msg.Payload.Service)

			value, err := msg.Payload.GetFloatValue()
			require.NoError(t, err)
			require.Equal(t, float64(100), value)

			return
		}
	}

}

func sampleAddress(t *testing.T, add int64) *fimpgo.Address {
	t.Helper()

	a, err := fimpgo.NewAddressFromString(fmt.Sprintf("pt:j1/mt:evt/rt:dev/rn:iqcontrols/ad:1/sv:sensor_lumin/ad:%d_0", add))
	require.NoError(t, err)

	return a
}
