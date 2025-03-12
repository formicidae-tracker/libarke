package main

import "github.com/formicidae-tracker/libarke/src-go/arke"

func init() {
	zeusCommand := MustAddCommand(parser.Command, "zeus",
		"Zeus command group.",
		"A collection of commands that can be sent to zeus devices.",
		nodeID)

	MustAddCommand(zeusCommand, "setPoint",
		"Sets the set point of zeus devices",
		"Sets the set point of zeus devices. A setPoint consist of an humidity (%R.H.) a temperature (°C) and a wind level (0-255).",
		&ArkeCommand[*arke.ZeusSetPoint]{Args: &arke.ZeusSetPoint{}})

	zeusGetCommand := MustAddCommand(getCommand, "zeus",
		"Zeus request group",
		"A collection of requests to ask data from zeus devices",
		nodeID)

	MustAddCommand(zeusGetCommand, "setPoint",
		"Requests zeus set point",
		"Requests zeus set point. A set point consists of humidity (%R.H.), temperature (°C), and wind level (0-255)",
		&Request[*arke.ZeusSetPoint]{message: &arke.ZeusSetPoint{}})

	MustAddCommand(zeusGetCommand, "report",
		"Requests a zeus report",
		"Requests a zeus report. A zeus reports consists of the humidity (R.H.), and temperatures",
		&Request[*arke.ZeusReport]{message: &arke.ZeusReport{}})

	MustAddCommand(zeusGetCommand, "config",
		"Requests zeus config",
		"Requests zeus config. A zeus config are the PID constant for the humidity and temperature control.",
		&Request[*arke.ZeusConfig]{message: &arke.ZeusConfig{}})

	MustAddCommand(zeusGetCommand, "status",
		"Requests zeus status",
		"Requests zeus status. It consists of the climate control status (idle|active|issues), and the fan status (RPM and age)",
		&Request[*arke.ZeusStatus]{message: &arke.ZeusStatus{}})

	MustAddCommand(zeusGetCommand, "controlPoint",
		"Requests zeus control commands",
		"Requests zeus current control value. It consists of two 16 bytes word, giving the heat / cooling and humidification desired power",
		&Request[*arke.ZeusControlPoint]{message: &arke.ZeusControlPoint{}})

	MustAddCommand(zeusGetCommand, "deltas",
		"Requests zeus temperature deltas",
		"Requests zeus temperature deltas. It is 4 offsets in °C to add to each sensors",
		&Request[*arke.ZeusDeltaTemperature]{message: &arke.ZeusDeltaTemperature{}})

}
