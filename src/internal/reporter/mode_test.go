package reporter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/stretchr/testify/require"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
)

func TestMode_Report(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		msg      *fimpgo.Message
		expected string
		wantErr  bool
	}{
		{
			name:     "home mode",
			msg:      homeModeMessage(t, "home"),
			expected: "40",
		},
		{
			name:     "away mode",
			msg:      homeModeMessage(t, "away"),
			expected: "0",
		},
		{
			name:     "sleep mode",
			msg:      homeModeMessage(t, "sleep"),
			expected: "0",
		},
		{
			name:     "vacation mode",
			msg:      homeModeMessage(t, "vacation"),
			expected: "0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/cgx/eepl.cgx" {
					require.Equal(t, tt.expected, r.FormValue("EEV0002"))
				} else {
					require.Equal(t, "/html_old/cgx/ios.cgx", r.URL.Path)
					require.NotEmpty(t, r.FormValue("IOO0122"))
				}

				require.Equal(t, "POST", r.Method)
				require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			}))

			defer s.Close()

			m := reporter.NewMode(s.URL, http.DefaultClient)
			err := m.Report(tt.msg)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func homeModeMessage(t *testing.T, mode string) *fimpgo.Message {
	t.Helper()

	msg, err := fimpgo.NewMessageFromBytes([]byte(fmt.Sprintf(`{
	  "serv": "vinculum",
	  "type": "evt.pd7.notify",
	  "val_t": "object",
	  "val": {
		"cmd": "set",
		"component": "hub",
		"id": "mode",
		"param": {
		  "current": "%s",
		  "prev": "home"
		}
	  },
	  "props": {},
	  "tags": null,
	  "src": "-",
	  "ver": "1",
	  "uid": "e8efec52-cf47-422d-a070-a90ba636698b",
	  "topic": "pt:j1/mt:evt/rt:app/rn:vinculum/ad:1"
	}`, mode)))
	require.NoError(t, err)

	return &fimpgo.Message{
		Addr: &fimpgo.Address{
			PayloadType:     "j1",
			MsgType:         "evt",
			ResourceType:    "app",
			ResourceName:    "vinculum",
			ResourceAddress: "1",
		},
		Payload: msg,
	}
}
