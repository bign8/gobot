package serial

// TODO: verify any of this works

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

var _ gobot.Adaptor = (*pin)(nil)
var _ Connection = (*pin)(nil)

type pin struct {
	name        string
	read, write *gpio.DirectPinDriver
}

// NewPinConnection constructs a serial port adapter based on a pair of gpio pins
func NewPinConnection(name string, reader, writer *gpio.DirectPinDriver) Connection {
	return &pin{
		name:  name,
		read:  reader,
		write: writer,
	}
}

// Name returns the label for the Adaptor
func (p *pin) Name() string {
	return p.name
}

// Port returns the concatenation of the adaptors ports
func (p *pin) Port() string {
	return p.read.Pin() + " " + p.write.Pin()
}

// Connect does stuff (TODO)
func (p *pin) Connect() (errs []error) {
	return nil
}

// Finalize does stuff (TODO)
func (p *pin) Finalize() (errs []error) {
	return nil
}

// DigitalWrite writes a byte to the serial buffer
func (p *pin) DigitalWrite(bits byte) (err error) {
	return p.write.DigitalWrite(bits)
}

// DigitalRead reads a byte from the serial buffer
func (p *pin) DigitalRead() (val int, err error) {
	return p.read.DigitalRead()
}
