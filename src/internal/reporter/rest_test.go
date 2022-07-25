package reporter_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	. "github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter/mocks"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/tests"
)

func TestRest_Report(t *testing.T) {
	t.Parallel()

	pub := mocks.NewRestPublisher(t)
	mapper := mocks.NewDeviceMapper(t)

	mapper.On("Device", mock.Anything).Return("IOO005", "1000", nil)
	pub.On("Publish", "IOO005", "1000").Return(nil)

	r := NewRest(pub, mapper)

	err := r.Report(binSwitchCommand(5, true))
	assert.NoError(t, err)
}

func TestRest_ReportErrored(t *testing.T) {
	t.Parallel()

	mapper := mocks.NewDeviceMapper(t)

	mapper.On("Device", mock.Anything).Return("", "", errors.New("test error"))

	r := NewRest(nil, mapper)
	err := r.Report(binSwitchCommand(5, true))
	require.Error(t, err)
}

func TestRestPublisher_Publish(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		host    string
		address string
		value   string
		wantErr bool
	}{
		{
			name:    "publish value 123 to IOI00432",
			address: "IOI00432",
			value:   "123",
		},
		{
			name:    "bad mass host",
			host:    "http$://localhost",
			address: "IOI00432",
			value:   "123",
			wantErr: true,
		},
	}

	for _, ttt := range testCases {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "POST", r.Method)
				require.Equal(t, "/cgx/custom/all.cgx", r.URL.Path)
				require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
				require.Equal(t, tt.value, r.FormValue(tt.address))
			}))

			defer s.Close()

			host := tt.host
			if host == "" {
				host = s.URL
			}

			publisher := NewRestPublisher(host, http.DefaultClient)

			err := publisher.Publish(tt.address, tt.value)
			if tt.wantErr {
				require.Error(t, err)

				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRestPublisher_Publish_ErroredResponse(t *testing.T) {
	t.Parallel()

	d := mocks.NewDoer(t)
	d.On("Do", mock.Anything).Return(nil, errors.New("test error"))

	publisher := NewRestPublisher("http://localhost", d)

	err := publisher.Publish("IOI00432", "123")
	require.Error(t, err)
}

func binSwitchCommand(addr int, val bool) *fimpgo.Message {
	msg := fimpgo.NewBoolMessage("cmd.binary.set", "out_bin_switch", val, nil, nil, nil)

	return &fimpgo.Message{
		Addr:    tests.DeviceAddress(addr),
		Payload: msg,
	}
}
