package roomba

import (
	"fmt"
	"time"
)

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

// Start starts the OI. You must always send the Start command before sending
// any other commands to the OI.
func (r *Driver) Start() []error {
	return r.write(128)
}

// Baud sets the baud rate in bits per second (bps) at which OI commands and
// data are sent according to the baud code sent in the data byte. The default
// baud rate at power up is 115200 bps, but the starting baud rate can be
// changed to 19200 by holding down the Clean button while powering on Roomba
// until you hear a sequence of descending tones. Once the baud rate is changed,
// it persists until Roomba is power cycled by pressing the power button or
// removing the battery, or when the battery voltage falls below the minimum
// required for processor operation. You must wait 100ms after sending this
// command before sending additional commands at the new baud rate.
//
// baudCode is defined using the following table.
//
//  *-----------*-----------*
//  | Baud Code | Baud Rate |
//  *-----------*-----------*
//  |        0  |     300   |
//  |        1  |     600   |
//  |        2  |    1200   |
//  |        3  |    2400   |
//  |        4  |    4800   |
//  |        5  |    9600   |
//  |        6  |   14400   |
//  |        7  |   19200   |
//  |        8  |   28800   |
//  |        9  |   38400   |
//  |       10  |   57600   |
//  |       11  |  115200   |
//  *-----------*-----------*
//
func (r *Driver) Baud(baudCode uint8) []error {
	if baudCode > 11 {
		return push(nil, fmt.Errorf("Invalid Baud Code: %d", baudCode))
	}
	return r.write(129, baudCode)
}

// Note (130): The effect and usage of the Control command (130) are identical
// to the Safe command.

// Safe puts the OI into Safe mode, enabling user control of Roomba. It turns
// off all LEDs. The OI can be in Passive, Safe, or Full mode to accept this
// command. If a safety condition occurs (see above) Roomba reverts
// automatically to Passive mode.
func (r *Driver) Safe() []error {
	return r.write(131)
}

// Full gives you complete control over Roomba by putting the OI into Full mode,
// and turning off the cliff, wheel-drop and internal charger safety features.
// That is, in Full mode, Roomba executes any command that you send it, even if
// the internal charger is plugged in, or command triggers a cliff or wheel drop
// condition.
func (r *Driver) Full() []error {
	return r.write(132)
}

// Power powers down Roomba. The OI can be in Passive, Safe, or Full mode to
// accept this command.
func (r *Driver) Power() []error {
	return r.write(133)
}

// Spot starts the Spot cleaning mode.
func (r *Driver) Spot() []error {
	return r.write(134)
}

// Clean starts the default cleaning mode.
func (r *Driver) Clean() []error {
	return r.write(135)
}

// Max starts the Max cleaning mode.
//
// Note: some documentation show this as a DEMO command (RESEARCH)
func (r *Driver) Max() []error {
	return r.write(136)
}

// Drive controls Roomba’s drive wheels. It takes four data bytes, interpreted
// as two 16-bit signed values using two’s complement. The first two bytes
// specify the average velocity of the drive wheels in millimeters per second
// (mm/s), with the high byte being sent first. The next two bytes specify the
// radius in millimeters at which Roomba will turn. The longer radii make Roomba
// drive straighter, while the shorter radii make Roomba turn more. The radius
// is measured from the center of the turning circle to the center of Roomba. A
// Drive command with a positive velocity and a positive radius makes Roomba
// drive forward while turning toward the left. A negative radius makes Roomba
// turn toward the right. Special cases for the radius make Roomba turn in place
// or drive straight, as specified below. A negative velocity makes Roomba drive
// backward.
//
// Note: Internal and environmental restrictions may prevent Roomba from
// accurately carrying out some drive commands. For example, it may not be
// possible for Roomba to drive at full speed in an arc with a large radius of
// curvature.
//
// Radius Special Cases
// - Straight = -32768 or 32767 = hex 8000 or 7FFF
// - Turn in place clockwise = -1
// - Turn in place counter-clockwise = 1
//
func (r *Driver) Drive(velocity int16, radius int16) (errs []error) {
	if velocity < -500 || velocity > 500 {
		errs = append(errs, fmt.Errorf("Invalid Velocity: %d", velocity))
	}
	if (radius < -2000 || radius > 2000) && radius != -32768 && radius != 32767 {
		errs = append(errs, fmt.Errorf("Invalid Radius: %d", radius))
	}
	if errs != nil {
		return errs
	}
	return r.write(137,
		uint8((velocity>>8)&0xFF),
		uint8((velocity>>0)&0xFF),
		uint8((radius>>8)&0xFF),
		uint8((radius>>0)&0xFF),
	)
}

// Motors lets you control the forward and backward motion of Roomba’s main
// brush, side brush, and vacuum independently. Motor velocity cannot be
// controlled with this command, all motors will run at maximum speed when
// enabled. The main brush and side brush can be run in either direction. The
// vacuum only runs forward.
//
// Serial Sequence: [138] [Motors]
//
// Motors is set in the following method
//  *-------*-----*-----*-----*-----*-----*-----*-----*-----*
//  |  Bit  |  7  |  6  |  5  |  4  |  3  |  2  |  1  |  0  |
//  *-------*-----*-----*-----*-----*-----*-----*-----*-----*
//  | Value | n/a | n/a | n/a | MBD | SBC |  MB |  V  |  SB |
//  *-------*-----*-----*-----*-----*-----*-----*-----*-----*
//
//  - MBD - Main Brush Direction
//  - SBC - Side Brush clockwise
//  - MB - Main Brush
//  - V - Vacuum
//  - SB - Side Brush
func (r *Driver) Motors(mainBrushDirection bool, sideBrushClockwise bool, mainBrush bool, vacuum bool, sideBrush bool) []error {
	var motorbits uint8 = 0x00
	args := []bool{sideBrush, vacuum, mainBrush, sideBrushClockwise, mainBrushDirection}
	for i, v := range args {
		if v {
			motorbits |= 0x01 << (uint)(i)
		}
	}
	return r.write(138, motorbits)
}

// LEDs ontrols the LEDs common to all models of Roomba 500. The Clean/Power LED
// is specified by two data bytes: one for the color and the other for the
// intensity.
//
// Serial Sequence: [139] [LED Bits] [Clean/Power Color] [Clean/Power Intensity]
//
// LED Bits is defined as follows
//  *-------*-----*-----*-----*-----*-------------*------*------*--------*
//  |  Bit  |  7  |  6  |  5  |  4  |      3      |   2  |   1  |    0   |
//  *-------*-----*-----*-----*-----*-------------*------*------*--------*
//  | Value | n/a | n/a | n/a | n/a | Check Robot | Dock | Spot | Debris |
//  *-------*-----*-----*-----*-----*-------------*------*------*--------*
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

// Song lets you specify up to four songs to the OI that you can play at a later
// time. Each song is associated with a song number. The Play command uses the
// song number to identify your song selection. Each song can contain up to
// sixteen notes. Each note is associated with a note number that uses MIDI note
// definitions and a duration that is specified in fractions of a second. The
// number of data bytes varies, depending on the length of the song specified. A
// one note song is specified by four data bytes. For each additional note
// within a song, add two data bytes.
func (r *Driver) Song(songNumber uint8, notes []Note) []error {
	data := []uint8{songNumber, (uint8)(len(notes))} // TODO: pre-allocate array

	for _, note := range notes {
		data = append(data, note.Number, note.Duration)
	}
	return r.write(140, data...)
}

// Play lets you select a song to play from the songs added to Roomba using the
// Song command. You must add one or more songs to Roomba using the Song command
// in order for the Play command to work.
func (r *Driver) Play(songNumber uint8) []error {
	// TODO: validate input
	return r.write(141, songNumber)
}

// Sensors requests the OI to send a packet of sensor data bytes. There are 58
// different sensor data packets. Each provides a value of a specific sensor or
// group of sensors.
func (r *Driver) Sensors(packetID uint8) []error {
	// TODO: validate input
	return r.write(142, packetID)
}

// SeekDock sends Roomba to the dock.
func (r *Driver) SeekDock() []error {
	return r.write(143)
}

// PWMMotors lets you control the speed of Roomba’s main brush, side brush, and
// vacuum independently. With each data byte, you specify the duty cycle for the
// low side driver (max 128). For example, if you want to control a motor with
// 25% of battery voltage, choose a duty cycle of 128 * 25% = 32. The main brush
// and side brush can be run in either direction. The vacuum only runs forward.
// Positive speeds turn the motor in its default (cleaning) direction. Default
// direction for the side brush is counterclockwise. Default direction for the
// main brush/flapper is inward.
//
// Serial sequence: [144] [Main Brush PWM] [Side Bruch PWM] [Vaccuum PWM]
func (r *Driver) PWMMotors(mainBrushPWM int8, sideBrushPWM int8, vacuumPWM uint8) []error {
	if vacuumPWM > 127 {
		return []error{fmt.Errorf("Invalid vaccumPWM: %d", vacuumPWM)}
	}
	return r.write(144, uint8(mainBrushPWM), uint8(sideBrushPWM), uint8(vacuumPWM))
}

// DriveDirect lets you control the forward and backward motion of Roomba’s
// drive wheels independently. It takes four data bytes, which are interpreted
// as two 16-bit signed values using two’s complement. The first two bytes
// specify the velocity of the right wheel in millimeters per second (mm/s),
// with the high byte sent first. The next two bytes specify the velocity of the
// left wheel, in the same format. A positive velocity makes that wheel drive
// forward, while a negative velocity makes it drive backward.
func (r *Driver) DriveDirect(rightVelocity int16, leftVelocity int16) (errs []error) {
	if rightVelocity < -500 || rightVelocity > 500 {
		errs = append(errs, fmt.Errorf("Invalid rightVelocity: %d", rightVelocity))
	}
	if leftVelocity < -500 || leftVelocity > 500 {
		errs = append(errs, fmt.Errorf("Invalid leftVelocity: %d", leftVelocity))
	}
	if errs != nil {
		return errs
	}
	return r.write(145,
		uint8((rightVelocity>>8)&0xFF),
		uint8((rightVelocity>>0)&0xFF),
		uint8((leftVelocity>>8)&0xFF),
		uint8((leftVelocity>>0)&0xFF),
	)
}

// DrivePWM lets you control the raw forward and backward motion of Roomba’s
// drive wheels independently. It takes four data bytes, which are interpreted
// as two 16-bit signed values using two’s complement. The first two bytes
// specify the PWM of the right wheel, with the high byte sent first. The next
// two bytes specify the PWM of the left wheel, in the same format. A positive
// PWM makes that wheel drive forward, while a negative PWM makes it drive
// backward.
func (r *Driver) DrivePWM(rightPWM int16, leftPWM int16) (errs []error) {
	if rightPWM < -255 || rightPWM > 255 {
		errs = append(errs, fmt.Errorf("Invalid rightPWM: %d", rightPWM))
	}
	if leftPWM < -255 || rightPWM > 255 {
		errs = append(errs, fmt.Errorf("Invalid leftPWM: %d", leftPWM))
	}
	if errs != nil {
		return errs
	}
	return r.write(146,
		uint8((rightPWM>>8)&0xFF),
		uint8((rightPWM>>0)&0xFF),
		uint8((leftPWM>>8)&0xFF),
		uint8((leftPWM>>0)&0xFF),
	)
}

// DigitalOutput controls the state of the 3 digital output pins on the 25 pin
// Cargo Bay Connector. The digital outputs can provide up to 20 mA of current.
//
// Warning: When the Robot is switched ON, the Digital Outputs are High for the
// first 3 seconds during the initialization of the bootloader
func (r *Driver) DigitalOutput(bit1, bit2, bit3 bool) []error {
	var bits uint8 = 0x00
	args := []bool{bit1, bit2, bit3}
	for i, v := range args {
		if v {
			bits |= 0x01 << uint(i)
		}
	}
	return r.write(147, bits)
}

// Stream starts a stream of data packets. The list of packets requested is sent
// every 15 ms, which is the rate Roomba uses to update data.
func (r *Driver) Stream(packetIDs []uint8) []error {
	// TODO: validate packetIDs
	packetIDs = append([]uint8{uint8(len(packetIDs))}, packetIDs...)
	return r.write(148, packetIDs...)
}

// QueryList lets you ask for a list of sensor packets. The result is returned
// once, as in the Sensors command. The robot returns the packets in the order
// you specify.
func (r *Driver) QueryList(packetIDs []uint8) []error {
	// TODO: validate packetIDs
	return r.write(149, packetIDs...)
}

// PauseStream lets you stop the steam without clearing the list of requested
// packets.
func (r *Driver) PauseStream() []error {
	return r.write(150, 0x00)
}

// ResumeStream lets you restart the steam without clearing the list of
// requested packets.
func (r *Driver) ResumeStream() []error {
	return r.write(150, 0x01)
}

// SendIR sends the requested byte out of low side driver 1 (pin 23 on the Cargo
// Bay Connector), using the format expected by iRobot Create’s IR receiver. You
// must use a preload resistor (suggested value: 100 ohms) in parallel with the
// IR LED and its resistor in order turn it on.
func (r *Driver) SendIR(msg byte) []error {
	return r.write(151, msg)
}

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

// Schedule sends Roomba a new schedule.
//
// days[0] = sun ... days[6] = sat
//
// Uses 15 data bytes:
//  [Days] [Sun Hour] [Sun Minute] [Mon Hour] [Mon Minute]
//         [Tue Hour] [Tue Minute] [Wed Hour] [Wed Minute]
//         [Thu Hour] [Thu Minute] [Fri Hour] [Fri Minute]
//         [Sat Hour] [Sat Minute]
//
// Days bit is setup in the following format
//  *-------*-----*-----*-----*-----*-----*-----*-----*-----*
//  |  Bit  |  7  |  6  |  5  |  4  |  3  |  2  |  1  |  0  |
//  *-------*-----*-----*-----*-----*-----*-----*-----*-----*
//  | Value | n/a | Sat | Fri | Thu | Wed | Tue | Mon | Sun |
//  *-------*-----*-----*-----*-----*-----*-----*-----*-----*
//
func (r *Driver) Schedule(days [7]time.Time) []error {
	var daybits uint8 = 0x00
	daybytes := []uint8{} // TODO: pre-allocate this space
	for i, day := range days {
		if !day.IsZero() {
			daybits |= 0x01 << uint8(i)
		}
		daybytes = append(daybytes, uint8(day.Hour()), uint8(day.Minute()))
	}
	return r.write(167, append([]uint8{daybits}, daybytes...)...)
}

// DisableSchedule disables scheduled cleaning, by sending all 0s.
func (r *Driver) DisableSchedule() []error {
	var days [7]time.Time
	return r.Schedule(days)
}

// SetDateTime sets Roomba’s clock.
func (r *Driver) SetDateTime(dateTime time.Time) []error {
	return r.write(168,
		uint8(dateTime.Weekday()),
		uint8(dateTime.Hour()),
		uint8(dateTime.Minute()),
	)
}

func (r *Driver) write(command uint8, data ...uint8) (errs []error) {
	d := append([]uint8{command}, data...)
	_, e := r.adap.sp.Write(d)
	return push(&errs, e)
}