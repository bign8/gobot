package roomba

import "github.com/hybridgroup/gobot"

// Driver describes a roomba driver for gobot
type Driver struct {
	name string
	adap *Adaptor
}

// NewDriver constructs a new gobot driver for a roomba
func NewDriver(a *Adaptor, name string) *Driver {
	return &Driver{
		name: name,
		adap: a,
	}
}

// Name sreturns the label for the Driver
func (r *Driver) Name() string {
	return r.name
}

// Halt terminates the Driver
func (r *Driver) Halt() []error {
	return r.Safe()
}

// Connection returns the Connection associated with the Driver
func (r *Driver) Connection() gobot.Connection {
	return r.adap
}

// Note: Start is in commands.go
