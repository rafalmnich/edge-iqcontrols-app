package decoder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/rafalmnich/edge-iqcontrols-app/internal/decoder"
)

func TestAddressAndValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []byte
		wantAddr int64
		wantVal  int64
	}{
		{
			name:     "address 130, value is 110",
			input:    []byte{170, 170, 170, 1, 16, 10, 130, 2, 110, 0},
			wantVal:  110,
			wantAddr: 130,
		},
		{
			name:     "value > 255",
			input:    []byte{170, 170, 170, 1, 16, 10, 130, 2, 2, 1},
			wantVal:  258,
			wantAddr: 130,
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			addr, val, _ := AddressAndValue(tt.input)
			assert.Equal(t, tt.wantAddr, addr)
			assert.Equal(t, tt.wantVal, val)

		})
	}
}
