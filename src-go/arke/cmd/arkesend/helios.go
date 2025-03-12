package main

import "github.com/formicidae-tracker/libarke/src-go/arke"

func init() {
	heliosCommand := MustAddCommand(parser.Command, "helios",
		"helios command group", "sends a command to helios devices",
		nodeID)
	MustAddCommand(heliosCommand, "setPoint",
		"sends wanted helios set point",
		"Sends helios set point",
		&ArkeCommand[*arke.HeliosSetPoint]{})
	MustAddCommand(heliosCommand, "pulse", "sets pulse mode for helios",
		"Sets pulse mode for helios. A duration of 0s disable pulse mode",
		&ArkeCommand[*arke.HeliosPulseMode]{})

	MustAddCommand(heliosCommand, "trigger", "sets trigger mode for helios",
		"Sets trigger mode for helios. A period of 0s disable trigger generation and enables external triggers",
		&ArkeCommand[*arke.HeliosTriggerMode]{})

	heliosGetCommand := MustAddCommand(getCommand, "helios",
		"helios command group", "request data from helios devices",
		nodeID)

	MustAddCommand(heliosGetCommand, "setPoint",
		"requests helios set point",
		"Requests helios set point",
		&Request[*arke.HeliosSetPoint]{})

	MustAddCommand(heliosGetCommand, "pulse", "request helios pulse mode",
		"Requests helios pulse mode. A period of 0s indicates no pulse",
		&Request[*arke.HeliosPulseMode]{})

	MustAddCommand(heliosGetCommand, "trigger", "request helios trigger mode",
		"Requests helios trigger mode. A period of 0s indicates no pulse",
		&Request[*arke.HeliosTriggerMode]{})

}
