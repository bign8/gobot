package roomba

import "github.com/hybridgroup/gobot"

// Driver describes a roomba driver for gobot
type Driver struct {
	gobot.Eventer
	gobot.Commander
	name string
	adap *Adaptor
}

// NewDriver constructs a new gobot driver for a roomba
func NewDriver(a *Adaptor, name string) *Driver {
	d := &Driver{
		Eventer:   gobot.NewEventer(),
		Commander: gobot.NewCommander(),
		name:      name,
		adap:      a,
	}

	// TODO: add all possible events
	// TODO: add all pottible commands

	return d
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
