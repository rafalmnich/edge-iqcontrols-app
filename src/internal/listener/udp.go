package listener

import (
	"io"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
)

type (
	// Connection represents
	Connection interface {
		io.Closer
		ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	}

	// Decoder represents a byte decoder for address and value.
	Decoder interface {
		Decode(b []byte) (add int64, val int64, err error)
	}

	// DecoderFunc represents a function that decodes bytes into address and value.
	DecoderFunc func(b []byte) (add int64, val int64, err error)
)

// Decode decodes bytes into address and value.
func (d DecoderFunc) Decode(b []byte) (add int64, val int64, err error) {
	return d(b)
}

// Listener represents a Listener listener.
type Listener interface {
	Start()
	Stop()
}

type udp struct {
	conn     Connection
	decoder  Decoder
	reporter reporter.Reporter

	stop chan struct{}
	done chan struct{}
}

// NewUDP creates a new UDP listener.
func NewUDP(conn Connection, decoder Decoder, reporter reporter.Reporter) Listener {
	return &udp{conn: conn, decoder: decoder, reporter: reporter, stop: make(chan struct{}), done: make(chan struct{})}
}

// Start starts the UDP listener.
func (u *udp) Start() {
	buf := make([]byte, 10)

	data := make(chan []byte, 1000)

	go u.listen(data, buf)

	for {
		select {
		case buf := <-data:
			add, val, err := u.decoder.Decode(buf)
			if err != nil {
				log.Errorf("error decoding: %s", err)

				continue
			}

			err = u.reporter.Report(add, val)
			if err != nil {
				log.Errorf("error reporting: %s", err)
			}
		case <-u.stop:
			close(u.done)

			return
		}
	}
}

// Stop stops the UDP listener.
func (u *udp) Stop() {
	close(u.stop)

	log.Info("waiting for Listener listener to stop")

	<-u.done

	log.Info("Listener listener stopped")

	err := u.conn.Close()
	if err != nil {
		log.Errorf("error closing Listener connection: %s", err)
	}
}

func (u *udp) listen(data chan []byte, buf []byte) {
	for {
		if _, _, err := u.conn.ReadFromUDP(buf); err != nil {
			log.Errorf("error reading from Listener: %s", err)

			return
		}

		select {
		case data <- buf:
			continue
		case <-u.stop:
			return
		}
	}
}
