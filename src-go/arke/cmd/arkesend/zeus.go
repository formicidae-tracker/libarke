package main

import "github.com/formicidae-tracker/libarke/src-go/arke"

func init() {
	zeusCommand := MustAddCommand(parser.Command, "zeus",
		"zeus command group", "zeus command group",
		nodeID)

	MustAddCommand(zeusCommand, "setPoint",
		"sets the set point of zeus devices",
		"Sets the set point of zeus devices",
		&ArkeCommand[*arke.ZeusSetPoint]{Args: &arke.ZeusSetPoint{}})

	zeusGetCommand := MustAddCommand(getCommand, "zeus",
		"zeus command group",
		"Zeus command group",
		nodeID)

	MustAddCommand(zeusGetCommand, "setPoint",
		"requests zeus set point",
		"Requests zeus set point",
		&Request[*arke.ZeusSetPoint]{message: &arke.ZeusSetPoint{}})
}
