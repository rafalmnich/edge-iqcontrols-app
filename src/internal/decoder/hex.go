package decoder

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// AddressAndValue converts hex parts of the input integer values of address and value
func AddressAndValue(b []byte) (addr int64, val int64, err error) {
	h := fmt.Sprintf("%x", b)
	h = h[len(h)-8:]

	addr, err = strconv.ParseInt(h[:2], 16, 64)
	if err != nil {
		return
	}

	multiplier, err := strconv.ParseInt(h[6:], 16, 64)
	if err != nil {
		return
	}

	val, err = strconv.ParseInt(h[4:6], 16, 64)

	val = val + multiplier*256

	log.Debugf("address: %d, value: %d", addr, val)

	return
}
