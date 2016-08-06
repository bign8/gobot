package serial

import (
	"io"

	"github.com/hybridgroup/gobot"
	serial "github.com/tarm/goserial"
)

var _ gobot.Adaptor = (*socket)(nil)
var _ Connection = (*socket)(nil)

type socket struct {
	name   string
	config *serial.Config
	sp     io.ReadWriteCloser
}

// NewSocketConnection constructs a new serial port adapter
func NewSocketConnection(name, port string, baud int) Connection {
	return &socket{
		name: name,
		config: &serial.Config{
			Name: port,
			Baud: baud,
		},
	}
}

// Name returns the label for the Adaptor
func (a *socket) Name() string {
	return a.name
}

// Port returns the adaptor's port
func (a *socket) Port() string {
	return a.config.Name
}

// Connect initates the Adaptor
func (a *socket) Connect() (errs []error) {
	if a.sp != nil {
		errs = a.Finalize()
	}
	if errs == nil {
		var err error
		a.sp, err = serial.OpenPort(a.config)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Finalize terminates the Adaptor
func (a *socket) Finalize() (errs []error) {
	if a.sp != nil {
		err := a.sp.Close()
		a.sp = nil
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// DigitalWrite writes a byte to the serial buffer
func (a *socket) DigitalWrite(bits byte) (err error) {
	if a.sp == nil {
		return ErrNotInitialized
	}
	_, err = a.sp.Write([]byte{bits})
	return err
}

// DigitalRead reads a byte from the serial buffer
func (a *socket) DigitalRead() (val int, err error) {
	if a.sp == nil {
		return 0, ErrNotInitialized
	}
	// TODO: verify this works
	bits := make([]byte, 1)
	_, err = a.sp.Read(bits)
	return int(bits[0]), err
}
