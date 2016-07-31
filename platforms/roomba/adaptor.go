package roomba

import (
	"io"

	"github.com/tarm/goserial"
)

// Adaptor is a roomba adapter for gobot (Implements Adaptor and Porter interfaces)
type Adaptor struct {
	name string
	port string
	sp   io.ReadWriteCloser
}

// NewAdaptor constructs a new roomba adapter
func NewAdaptor(name string, port string) *Adaptor {
	return &Adaptor{
		name: name,
		port: port,
	}
}

// Name returns the label for the Adaptor
func (r *Adaptor) Name() string {
	return r.name
}

// Port returns the adaptor's port
func (r *Adaptor) Port() string {
	return r.port
}

// Connect initiates the Adaptor
func (r *Adaptor) Connect() (errs []error) {
	if r.sp != nil {
		errs = r.Finalize()
	}
	if errs == nil {
		var e error
		r.sp, e = serial.OpenPort(&serial.Config{Name: r.port, Baud: 115200})
		push(&errs, e)
	}
	return errs
}

// Finalize terminates the Adaptor
func (r *Adaptor) Finalize() (errs []error) {
	if r.sp != nil {
		push(&errs, r.sp.Close())
		r.sp = nil
	}
	return errs
}

func push(e *[]error, err error) []error {
	if err != nil {
		*e = append(*e, err)
	}
	return *e
}
