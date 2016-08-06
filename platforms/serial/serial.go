package serial

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

// Connection is a serial connection that can be communicated with
// NOTE: gpio.DirectPinDriver also implements this interface
// TODO: make a helper that accepts 2 pins for bi-directional communication
type Connection interface {
	gobot.Adaptor
	DigitalRead() (val int, err error)
	DigitalWrite(level byte) (err error)
}

// ErrNotInitialized means you should probably connect to the port!
var ErrNotInitialized = errors.New("Not Initialized")
