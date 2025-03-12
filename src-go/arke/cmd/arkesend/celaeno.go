package main

import (
	"github.com/formicidae-tracker/libarke/src-go/arke"
)

func init() {
	celaenoCommand := MustAddCommand(parser.Command, "celaeno",
		"Celaeno command group",
		"A collection of commands that can be sent to celaeno devices",
		nodeID)
	MustAddCommand(celaenoCommand,
		"setPoint",
		"Sends celaeno set point",
		"Sends celaeno set point. A set point is just a 0-255 int that sets the desired humidification power.",
		&ArkeCommand[*arke.CelaenoSetPoint]{Args: &arke.CelaenoSetPoint{}})

	MustAddCommand(celaenoCommand, "config",
		"Sends celaeno config",
		"Sends celaeno config. This config defines the ramp up time, the ramp down time, the minimum on duration, and the debounce time for the sensor level. these durations should not exceed ~65s",
		&ArkeCommand[*arke.CelaenoConfig]{Args: &arke.CelaenoConfig{}})

	getCelaenoCommand := MustAddCommand(getCommand, "celaeno",
		"Celaeno request group",
		"A collection of request to ask data from celaeno devices",
		nodeID)

	MustAddCommand(getCelaenoCommand,
		"setPoint",
		"Requests celaeno set point",
		"Request celaeno set point. It is a byte indicating the current desired level of humidification.",
		&Request[*arke.CelaenoSetPoint]{message: &arke.CelaenoSetPoint{}})

	MustAddCommand(getCelaenoCommand,
		"status",
		"Requests celaeno status",
		"Requests celaeno status. It consists of the water level and the fan status.",
		&Request[*arke.CelaenoStatus]{message: &arke.CelaenoStatus{}})

	MustAddCommand(getCelaenoCommand, "config",
		"Requests celaeno config",
		"Requests celaeno config. It consists of the ramp up time, ramp down time, minimum on time, and debounce time.",
		&Request[*arke.CelaenoConfig]{message: &arke.CelaenoConfig{}})

}
