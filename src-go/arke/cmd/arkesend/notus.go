package main

import (
	"github.com/formicidae-tracker/libarke/src-go/arke"
)

func init() {
	notusCommand := MustAddCommand(parser.Command, "notus",
		"Notus command group",
		"A collection of commands that can be sent to notus devices.",
		nodeID)

	MustAddCommand(notusCommand, "setPoint",
		"Sends notus set point",
		"Sends notus set point. It consist of a byte representing the desired heating power.",
		&ArkeCommand[*arke.NotusSetPoint]{Args: &arke.NotusSetPoint{}})

	MustAddCommand(notusCommand, "config",
		"Sends notus config",
		"Sends notus config. It consists of a ramp-down time, a minimum fan level (byte) when on, and the maximum allowed heating power (byte).",
		&ArkeCommand[*arke.NotusConfig]{Args: &arke.NotusConfig{}})

	getNotusCommand := MustAddCommand(getCommand, "notus",
		"Notus request group",
		"A collection of request to poll data from notus devices.",
		nodeID)

	getNotusCommand.AddCommand("setPoint",
		"Requests notus set point",
		"Requests notus set point, i.e. a byte representing the current heating power.",
		&Request[*arke.NotusSetPoint]{message: &arke.NotusSetPoint{}})

	getNotusCommand.AddCommand("config",
		"Requests notus config",
		"Requests notus config. It consists of a ramp down duration, and the minimum fan level and the maximum heating power. Both value are bytes.",
		&Request[*arke.NotusConfig]{message: &arke.NotusConfig{}})

}
