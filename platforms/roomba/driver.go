package roomba

import "github.com/hybridgroup/gobot"

// Driver describes a roomba driver for gobot
type Driver struct {
	name string
	a    *Adaptor
}

// NewDriver constructs a new gobot driver for a roomba
func NewDriver(a *Adaptor, name string) *Driver {
	return &Driver{
		name: name,
		a:    a,
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
	return r.a
}

//
// The following functions are derived from the iRobot command manual
//    Create: http://www.irobot.com/filelibrary/pdfs/hrd/create/Create%20Open%20Interface_v2.pdf
//    Create2: https://cdn-shop.adafruit.com/datasheets/create_2_Open_Interface_Spec.pdf
//
//   Note: There are differeing commands between the Create1 and Create2.
//   TODO: Enforce the different commands based on the Create version
//
// Any TODO's with a name followed by a number in params (ex: TODO: control 130)
//   comes directly from the iRobot Create Open Interface PDF (mentioned above)
//
// Any TODO's with the label: 'from previous implementation' come from the last
//   developmental library found: https://github.com/karota-project/gobot-roomba/tree/develop
//   NOTE: these commands are missing from the documentation and sould be tested thouroughly
//

// Start starts the OI. You must always send the Start command before sending any other commands to the OI.
func (r *Driver) Start() []error {
	return r.write(128)
}

// Baud sets the baud rate in bits per second (bps) at which OI commands and data are sent according to the baud code sent in the data byte. The default baud rate at power up is 115200 bps, but the starting baud rate can be changed to 19200 by holding down the Clean button while powering on Roomba until you hear a sequence of descending tones. Once the baud rate is changed, it persists until Roomba is power cycled by pressing the power button or removing the battery, or when the battery voltage falls below the minimum required for processor operation. You must wait 100ms after sending this command before sending additional commands at the new baud rate.
func (r *Driver) Baud(rate uint8) []error {
	// TODO: validate input
	return r.write(129, rate)
}

// TODO: control (130)

// Safe puts the OI into Safe mode, enabling user control of Roomba. It turns off all LEDs. The OI can be in Passive, Safe, or Full mode to accept this command. If a safety condition occurs (see above) Roomba reverts automatically to Passive mode.
func (r *Driver) Safe() []error {
	return r.write(131)
}

// Full gives you complete control over Roomba by putting the OI into Full mode, and turning off the cliff, wheel-drop and internal charger safety features. That is, in Full mode, Roomba executes any command that you send it, even if the internal charger is plugged in, or command triggers a cliff or wheel drop condition.
func (r *Driver) Full() []error {
	return r.write(132)
}

// TODO: from previous implementation (not in documentation)
// POWER powers down Roomba. The OI can be in Passive, Safe, or Full mode to accept this command.
// POWER command = 133
// func (r *Driver) Power() {
// 	r.sender(COMMAND_POWER, []uint8{})
// }

// Spot starts the Spot cleaning mode.
func (r *Driver) Spot() []error {
	return r.write(134)
}

// Clean starts the default cleaning mode.
func (r *Driver) Clean() []error {
	return r.write(135)
}

// Demo allows you to run one of the preset demo programs.
func (r *Driver) Demo(num uint8) []error {
	// TODO: validate input
	return r.write(136, num)
}

// Drive controls Roomba’s drive wheels. It takes four data bytes, interpreted as two 16-bit signed values using two’s complement. The first two bytes specify the average velocity of the drive wheels in millimeters per second (mm/s), with the high byte being sent first. The next two bytes specify the radius in millimeters at which Roomba will turn. The longer radii make Roomba drive straighter, while the shorter radii make Roomba turn more. The radius is measured from the center of the turning circle to the center of Roomba. A Drive command with a positive velocity and a positive radius makes Roomba drive forward while turning toward the left. A negative radius makes Roomba turn toward the right. Special cases for the radius make Roomba turn in place or drive straight, as specified below. A negative velocity makes Roomba drive backward.
func (r *Driver) Drive(velocity int16, radius int16) []error {
	// TODO: validate input
	return r.write(137,
		uint8((velocity>>8)&0xFF),
		uint8((velocity>>0)&0xFF),
		uint8((radius>>8)&0xFF),
		uint8((radius>>0)&0xFF),
	)
}

// Motors lets you control the forward and backward motion of Roomba’s main brush, side brush, and vacuum independently. Motor velocity cannot be controlled with this command, all motors will run at maximum speed when enabled. The main brush and side brush can be run in either direction. The vacuum only runs forward.
func (r *Driver) Motors(mainBrushDirection bool, sideBrushClockwise bool, mainBrush bool, vacuum bool, sideBrush bool) []error {
	var motorbits uint8 = 0x00
	// TODO: documentation references 7 output bits, verify this
	args := []bool{sideBrush, vacuum, mainBrush, sideBrushClockwise, mainBrushDirection}
	for i, v := range args {
		if v {
			motorbits |= 0x01 << (uint)(i)
		}
	}
	return r.write(138, motorbits)
}

// LEDs ontrols the LEDs common to all models of Roomba 500. The Clean/Power LED is specified by two data bytes: one for the color and the other for the intensity.
func (r *Driver) LEDs(checkRobot bool, dock bool, spot bool, debris bool, cleanPowerColor uint8, cleanPowerIntensity uint8) []error {
	var ledbits uint8 = 0x00
	args := []bool{debris, spot, dock, checkRobot}
	for i, v := range args {
		if v {
			ledbits |= 0x01 << uint(i)
		}
	}
	return r.write(139, ledbits, cleanPowerColor, cleanPowerIntensity)
}

// Song lets you specify up to four songs to the OI that you can play at a later time. Each song is associated with a song number. The Play command uses the song number to identify your song selection. Each song can contain up to sixteen notes. Each note is associated with a note number that uses MIDI note definitions and a duration that is specified in fractions of a second. The number of data bytes varies, depending on the length of the song specified. A one note song is specified by four data bytes. For each additional note within a song, add two data bytes.
func (r *Driver) Song(songNumber uint8, notes []Note) []error {
	data := []uint8{songNumber, (uint8)(len(notes))} // TODO: pre-allocate array

	for _, note := range notes {
		data = append(data, note.Number, note.Duration)
	}
	return r.write(140, data...)
}

// Play lets you select a song to play from the songs added to Roomba using the Song command. You must add one or more songs to Roomba using the Song command in order for the Play command to work.
func (r *Driver) Play(songNumber uint8) []error {
	// TODO: validate input
	return r.write(141, songNumber)
}

// Sensors requests the OI to send a packet of sensor data bytes. There are 58 different sensor data packets. Each provides a value of a specific sensor or group of sensors.
func (r *Driver) Sensors(packetID uint8) []error {
	// TODO: validate input
	return r.write(142, packetID)
}

// SeekDock sends Roomba to the dock.
func (r *Driver) SeekDock() []error {
	return r.write(143)
}

// PwmMotors lets you control the speed of Roomba’s main brush, side brush, and vacuum independently. With each data byte, you specify the duty cycle for the low side driver (max 128). For example, if you want to control a motor with 25% of battery voltage, choose a duty cycle of 128 * 25% = 32. The main brush and side brush can be run in either direction. The vacuum only runs forward. Positive speeds turn the motor in its default (cleaning) direction. Default direction for the side brush is counterclockwise. Default direction for the main brush/flapper is inward.
func (r *Driver) PwmMotors(mainBrushPwm int8, sideBrushPwm int8, vacuumPwm uint8) []error {
	// TODO: validate input
	return r.write(144, uint8(mainBrushPwm), uint8(sideBrushPwm), uint8(vacuumPwm))
}

// DriveDirect lets you control the forward and backward motion of Roomba’s drive wheels independently. It takes four data bytes, which are interpreted as two 16-bit signed values using two’s complement. The first two bytes specify the velocity of the right wheel in millimeters per second (mm/s), with the high byte sent first. The next two bytes specify the velocity of the left wheel, in the same format. A positive velocity makes that wheel drive forward, while a negative velocity makes it drive backward.
func (r *Driver) DriveDirect(rightVelocity int16, leftVelocity int16) []error {
	return r.write(145,
		uint8((rightVelocity>>8)&0xFF),
		uint8((rightVelocity>>0)&0xFF),
		uint8((leftVelocity>>8)&0xFF),
		uint8((leftVelocity>>0)&0xFF),
	)
}

// DrivePwm lets you control the raw forward and backward motion of Roomba’s drive wheels independently. It takes four data bytes, which are interpreted as two 16-bit signed values using two’s complement. The first two bytes specify the PWM of the right wheel, with the high byte sent first. The next two bytes specify the PWM of the left wheel, in the same format. A positive PWM makes that wheel drive forward, while a negative PWM makes it drive backward.
func (r *Driver) DrivePwm(rightPwm int16, leftPwm int16) []error {
	return r.write(146,
		uint8((rightPwm>>8)&0xFF),
		uint8((rightPwm>>0)&0xFF),
		uint8((leftPwm>>8)&0xFF),
		uint8((leftPwm>>0)&0xFF),
	)
}

// TODO: Digital Outputs (147)

// Stream starts a stream of data packets. The list of packets requested is sent every 15 ms, which is the rate Roomba uses to update data.
func (r *Driver) Stream(packetIDs []uint8) []error {
	// TODO: validate packetIDs
	packetIDs = append([]uint8{uint8(len(packetIDs))}, packetIDs...)
	return r.write(148, packetIDs...)
}

// QueryList lets you ask for a list of sensor packets. The result is returned once, as in the Sensors command. The robot returns the packets in the order you specify.
func (r *Driver) QueryList(packetIDs []uint8) []error {
	// TODO: validate packetIDs
	return r.write(149, packetIDs...)
}

// PauseStream lets you stop the steam without clearing the list of requested packets.
func (r *Driver) PauseStream() []error {
	return r.write(150, 0x00)
}

// ResumeStream lets you restart the steam without clearing the list of requested packets.
func (r *Driver) ResumeStream() []error {
	return r.write(150, 0x01)
}

// TODO: send IR (151)

// TODO: script (152)

// TODO: paly script (153)

// TODO: show script (154)

// TODO: wait time (155)

// TODO: wait distance (156)

// TODO: wait angle (157)

// TODO: wait event (158)

// TODO: 159, 160, 161 (not in any documentation)

// SchedulingLEDs controls the state of the scheduling LEDs present on the Roomba 560 and 570.
func (r *Driver) SchedulingLEDs(weekday uint8, schedule bool, clock bool, am bool, pm bool, colon bool) []error {
	// TODO: validate input
	var ledbits uint8 = 0x00
	args := []bool{colon, pm, am, clock, schedule}
	for i, v := range args {
		if v {
			ledbits |= 0x01 << uint(i)
		}
	}
	return r.write(162, weekday, ledbits)
}

// DigitLEDsRaw controls the four 7 segment displays on the Roomba 560 and 570.
func (r *Driver) DigitLEDsRaw(digit [4]uint8) []error {
	// TODO: validate input
	return r.write(163, digit[:4]...)
}

// DigitLEDsASCII controls the four 7 segment displays on the Roomba 560 and 570 using ASCII character codes. Because a 7 segment display is not sufficient to display alphabetic characters properly, all characters are an approximation, and not all ASCII codes are implemented.
func (r *Driver) DigitLEDsASCII(message string) []error {
	// TODO: verify input
	return r.write(164, []uint8((message + "    ")[:4])...)
}

// Buttons lets you push Roomba’s buttons. The buttons will automatically release after 1/6th of a second.
func (r *Driver) Buttons(clock bool, schedule bool, day bool, hour bool, dock bool, spot bool, clean bool) []error {
	var buttonbits uint8 = 0x00
	args := []bool{clean, spot, dock, hour, day, schedule, clock}
	for i, v := range args {
		if v {
			buttonbits |= 0x01 << uint(i)
		}
	}
	return r.write(165, buttonbits)
}

// TODO: 166 not in any documentation
// TODO: these remaining functions existed in past implementation but have no associated documentation.
// // days[0] = sun ... days[6] = sat
// func (r *Driver) Schedule(days [7]time.Time) {
// 	var daybits uint8 = 0x00
// 	daybytes := []uint8{} // TODO: pre-allocate this space
//
// 	for i, day := range days {
// 		if !day.IsZero() {
// 			daybits |= 0x01 << uint8(i)
// 		}
// 		daybytes = append(daybytes, uint8(day.Hour()), uint8(day.Minute()))
// 	}
//
// 	r.sender(167, append([]uint8{daybits}, daybytes...))
// }
//
// func (r *Driver) DisableSchedule() {
// 	r.sender(167, []uint8{
// 		0x00, 0x00, 0x00, 0x00, 0x00,
// 		0x00, 0x00, 0x00, 0x00, 0x00,
// 		0x00, 0x00, 0x00, 0x00, 0x00,
// 	})
// }
//
// func (r *Driver) SetDateTime(dateTime time.Time) {
// 	r.sender(168, []uint8{
// 		uint8(dateTime.Weekday()),
// 		uint8(dateTime.Hour()),
// 		uint8(dateTime.Minute()),
// 	})
// }

func (r *Driver) write(command uint8, data ...uint8) (errs []error) {
	d := append([]uint8{command}, data...)
	_, e := r.a.sp.Write(d)
	return push(&errs, e)
}
