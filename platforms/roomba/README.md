# Roomba
iRobot Create is an affordable STEM resource for educators, students and developers. Use Create to grasp the fundamentals of robotics, computer science and engineering.

For more information about the Roomba, checkout the [iRobot](http://www.irobot.com/) website.

This driver is based on [previous implementations](https://github.com/karota-project/gobot-roomba/tree/develop)

## Install
As documented, the Roomba communicates via an external Serial Port Mini-DIN Connector.  This port exposes a 5Volt level serial port which can be interfaced via GPIO, Serial, USB2Serial adaptor, etc.

### USB 2 Serial
This set of instructions is for the default iRobot USB 2 Serial adapter that comes packaged with the iRobot Create 2 system.  Other USB 2 Serial adapters will have to find driver install instructions elsewhere.

Simply Download drivers available at one of the following locations.

- [FTDChip](http://www.ftdichip.com/Drivers/VCP.htm)
- [PBX-Book](http://pbxbook.com/other/mac-tty.html)

Verify the serial port shows up using the following commands.
- Mac `ls ls /dev/tty.*`
- Win TODO
- Linux TODO

And if you want to test your stuff, checkout the [Python Tethered Driving](http://www.irobotweb.com/~/media/MainSite/PDFs/About/STEM/Create/Python_Tethered_Driving.pdf) demo on the iRobot website.


# Introduction to Roomba

### Charging Statuses

- Quick Amber Flash
  - If Roomba senses that its battery has been significantly discharged, Roomba will enter a special 16-hour refresh charge cycle. When the 16-hour refresh charge is initiated, the CLEAN light will quickly pulse red/amber. Do not interrupt the cycle once it has begun. The 16-hour refresh charge cycle is initiated by the robot and cannot be started manually.
  - Source : http://aurobo.com.au/article_info.php?articles_id=44

## How to Use

```go
package main

// TODO: put this together
```

# References
The adaptors were built off of several pieces of documentation.

- [Create Open Interface Spec](http://www.irobot.com/filelibrary/pdfs/hrd/create/Create%20Open%20Interface_v2.pdf)
- [Create2 Open Interface Spec](https://cdn-shop.adafruit.com/datasheets/create_2_Open_Interface_Spec.pdf)
- [Roomba 500 Open Interface Spec](http://irobot.lv/uploaded_files/File/iRobot_Roomba_500_Open_Interface_Spec.pdf)
