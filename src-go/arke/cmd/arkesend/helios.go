package main

import "github.com/formicidae-tracker/libarke/src-go/arke"

func init() {
	heliosCommand := MustAddCommand(parser.Command, "helios",
		"Helios command group",
		"A collection of commands that can be sent to helios devices",
		nodeID)
	MustAddCommand(heliosCommand, "setPoint",
		"Sends helios set point",
		"Sends helios set point. It consists of an amount of visible and UV light, both as a byte (0-255).",
		&ArkeCommand[*arke.HeliosSetPoint]{Args: &arke.HeliosSetPoint{}})
	MustAddCommand(heliosCommand, "pulse",
		"Sets pulse mode for helios",
		"Sets pulse mode for helios. current visible and uv light will pulse over the given period between zero and their assigned values. It requires a period, that should not exceed ~65s.",
		&ArkeCommand[*arke.HeliosPulseMode]{Args: &arke.HeliosPulseMode{}})

	MustAddCommand(heliosCommand, "trigger",
		"Sets the trigger mode for helios",
		"Sets the trigger mode for helios. It consist of a period and a pulse duration. If a period is given, helios will ignore external triggers and sends IR pulses according to the given period. period should not exceed ~6.5s and pulse should not exceed 3.5ms. A period of zero indicates external triggers",
		&ArkeCommand[*arke.HeliosTriggerMode]{Args: &arke.HeliosTriggerMode{}})

	heliosGetCommand := MustAddCommand(getCommand, "helios",
		"Helios request group",
		"A collection of request to ask data from helios devices.",
		nodeID)

	MustAddCommand(heliosGetCommand, "setPoint",
		"Requests helios set point",
		"Requests helios set point. It consist of an amount of visible and UV light",
		&Request[*arke.HeliosSetPoint]{message: &arke.HeliosSetPoint{}})

	MustAddCommand(heliosGetCommand, "pulse",
		"Requests helios pulse mode",
		"Requests helios pulse mode. A period of 0s indicates no pulse effect.",
		&Request[*arke.HeliosPulseMode]{message: &arke.HeliosPulseMode{}})

	MustAddCommand(heliosGetCommand, "trigger",
		"Requests helios trigger mode",
		"Requests helios trigger mode. A period of 0s indicates no self-generation of IR triggers.",
		&Request[*arke.HeliosTriggerMode]{message: &arke.HeliosTriggerMode{}})

}
