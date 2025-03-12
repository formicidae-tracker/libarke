package main

import (
	"github.com/formicidae-tracker/libarke/src-go/arke"
)

func init() {
	celaenoCommand := MustAddCommand(parser.Command, "celaeno",
		"celaeno command group",
		"sends a command for celaeno devices over the CANbus",
		nodeID)
	MustAddCommand(celaenoCommand,
		"setPoint",
		"sends celaeno set point",
		"Sends celaeno set point",
		&ArkeCommand[*arke.CelaenoSetPoint]{Args: &arke.CelaenoSetPoint{}})

	MustAddCommand(celaenoCommand, "config",
		"sends celaeno config",
		"Sends celaeno config",
		&ArkeCommand[*arke.CelaenoConfig]{Args: &arke.CelaenoConfig{}})

	getCelaenoCommand := MustAddCommand(getCommand, "celaeno",
		"celaeno command group",
		"sends a request for celaeno devices over the CANbus",
		nodeID)

	MustAddCommand(getCelaenoCommand,
		"setPoint",
		"requests celaeno set point",
		"Request celaeno set point",
		&Request[*arke.CelaenoSetPoint]{message: &arke.CelaenoSetPoint{}})

	MustAddCommand(getCelaenoCommand,
		"status",
		"requests celaeno status",
		"Requests celaeno status",
		&Request[*arke.CelaenoStatus]{message: &arke.CelaenoStatus{}})

	MustAddCommand(getCelaenoCommand, "config",
		"requests celaeno config",
		"Requests celaeno config",
		&Request[*arke.CelaenoConfig]{message: &arke.CelaenoConfig{}})

}
