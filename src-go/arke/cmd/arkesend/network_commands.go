package main

import (
	"time"

	"github.com/formicidae-tracker/libarke/src-go/arke"
)

type NetworkCommand struct {
	Class NodeClassName `long:"class" short:"c" default:"broadcast" description:"node class to target"`
}

var network = &NetworkCommand{}

type ResetCommand struct {
	ID uint8 `long:"ID" short:"I" default:"0" description:"ID to target, 0 to broadcast"`
}

func (cmd *ResetCommand) Execute(args []string) error {
	return opts.Send(arke.MakeResetRequest(network.Class.Class(), arke.NodeID(cmd.ID)))
}

type PingCommand struct {
}

func (cmd *PingCommand) Execute(args []string) error {
	return opts.Send(arke.MakePing(network.Class.Class()))
}

type HeartbeatCommand struct {
	Args struct {
		Period time.Duration `positional-arg-name:"period" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *HeartbeatCommand) Execute(args []string) error {
	return opts.Send(arke.MakeHeartBeatRequest(network.Class.Class(), cmd.Args.Period))
}

type ChangeIDCommand struct {
	Args struct {
		Old uint8 `positional-arg-name:"old" required:"yes"`
		New uint8 `positional-arg-name:"new" required:"yes"`
	} `positional-args:"yes"`
}

func (cmd *ChangeIDCommand) Execute([]string) error {
	return opts.Send(arke.MakeIDChangeRequest(network.Class.Class(),
		arke.NodeID(cmd.Args.Old), arke.NodeID(cmd.Args.New)))
}

func keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, len(m))
	i := 0
	for k := range m {
		res[i] = k
		i++
	}
	return res
}

func init() {
	networkCommand, _ := parser.AddCommand("network",
		"network command group",
		"sends a network command over the CANbus", network)
	networkCommand.FindOptionByLongName("class").Choices = keys(nodeClassName)
	networkCommand.AddCommand("reset",
		"sends a reset command",
		"sends a reset command to a given node", &ResetCommand{})
	networkCommand.AddCommand("ping",
		"ping a class of node",
		"Requests a single heartbeat command to a given class of nodes",
		&PingCommand{})
	networkCommand.AddCommand("heartbeat",
		"requests periodic heartbeats",
		"Requests a class of nodes to send periodic heartbeat",
		&HeartbeatCommand{})
}
