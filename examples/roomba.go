package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/roomba"
)

func main() {
	master := gobot.NewGobot()
	api.NewAPI(master).Start()

	adaptor := roomba.NewAdaptor("roomba-a01", "/dev/tty.usbserial-DA01NPOT")
	driver := roomba.NewDriver(adaptor, "roomba-d01")

	master.AddRobot(
		gobot.NewRobot(
			"roomba",
			[]gobot.Connection{adaptor},
			[]gobot.Device{driver},
			func() {
				fmt.Println("work")

				driver.Full()

				driver.DigitLEDsASCII("ABCD")

				time.Sleep(5 * time.Second)
				fmt.Println("finish")
				driver.Safe()
				driver.SeekDock()
			}))

	master.Start()
}
