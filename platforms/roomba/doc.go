/*
Package roomba contains the Gobot adaptor for iRobot Roombas.

For further information refer to raspi README:
https://github.com/hybridgroup/gobot/blob/master/platforms/roomba/README.md
*/
package roomba

func push(e *[]error, err error) []error {
	if e == nil {
		if err != nil {
			return []error{err}
		}
		return nil
	}

	if err != nil {
		*e = append(*e, err)
	}
	return *e
}
