package main

import (
	"fmt"
	"os"

	socketcan "github.com/atuleu/golang-socketcan"
	"github.com/formicidae-tracker/libarke/src-go/arke"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Interface    arke.CANInterfaceName `long:"interface" short:"i" default:"slcan0" description:"CAN interface to use"`
	HighPriority bool                  `long:"priority" short:"P"`
}

func (o *Options) buildStandardMessage(m arke.Message, n arke.NodeID, RTR bool) socketcan.CanFrame {
	mType := arke.StandardMessage
	if o.HighPriority == true {
		mType = arke.HighPriorityMessage
	}

	frame := socketcan.CanFrame{
		ID:       arke.MakeCANIDT(mType, m.MessageClassID(), n),
		RTR:      RTR,
		Extended: false,
		Data:     make([]byte, 8),
	}

	if RTR == false {
		size, _ := m.Marshal(frame.Data)
		frame.Dlc = byte(size)
	}

	return frame
}

type NodeIDGroup struct {
	ID arke.NodeID `short:"I" long:"ID" description:"ID to target" default:"0"`
}

var nodeID = &NodeIDGroup{}

type Request[M arke.Message] struct {
	message M
}

func (cmd *Request[M]) Execute([]string) error {
	return opts.Send(opts.buildStandardMessage(cmd.message, nodeID.ID, true))
}

type ArkeCommand[M arke.Message] struct {
	Args M `positional-args:"yes"`
}

func (cmd *ArkeCommand[M]) Execute([]string) error {
	return opts.Send(opts.buildStandardMessage(cmd.Args, nodeID.ID, false))
}

func (o *Options) Send(frame socketcan.CanFrame) error {
	intf, err := socketcan.NewRawInterface(string(o.Interface))
	if err != nil {
		return fmt.Errorf("opening CAN interface '%s': %s", o.Interface, err)
	}
	defer intf.Close()

	return intf.Send(frame)

}

var opts = &Options{}
var parser = flags.NewParser(opts, flags.Default)

func main() {
	if len(os.Getenv("GO_FLAGS_MANPAGE")) > 0 {
		parser.WriteManPage(os.Stdout)
		return
	}

	_, err := parser.Parse()
	if flags.WroteHelp(err) == true {
		return
	}
	if _, ok := err.(*flags.Error); ok == true {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	if err != nil {
		os.Exit(1)
	}
}
