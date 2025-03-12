package main

import (
	"github.com/formicidae-tracker/libarke/src-go/arke"
)

func init() {
	notusCommand := MustAddCommand(parser.Command, "notus",
		"notus command group",
		"sends a command for notus devices over the CANbus",
		nodeID)

	MustAddCommand(notusCommand, "setPoint", "sends notus set point",
		"Sends notus set point to a given node",
		&ArkeCommand[*arke.NotusSetPoint]{Args: &arke.NotusSetPoint{}})

	MustAddCommand(notusCommand, "config",
		"sends config to notus devices",
		"Sends config to notus devices over the CANbus",
		&ArkeCommand[*arke.NotusConfig]{Args: &arke.NotusConfig{}})

	getNotusCommand := MustAddCommand(getCommand, "notus",
		"notus command group",
		"sends a request for notus devices over the CANbus",
		nodeID)

	getNotusCommand.AddCommand("setPoint",
		"requests notus set point",
		"Requests notus set point",
		&Request[*arke.NotusSetPoint]{message: &arke.NotusSetPoint{}})

	getNotusCommand.AddCommand("config",
		"requests notus config",
		"Requests notus config",
		&Request[*arke.NotusConfig]{message: &arke.NotusConfig{}})

}
