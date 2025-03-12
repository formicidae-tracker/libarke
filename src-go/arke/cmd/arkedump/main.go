package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	socketcan "github.com/atuleu/golang-socketcan"
	"github.com/formicidae-tracker/libarke/src-go/arke"
	"github.com/jessevdk/go-flags"
	"golang.org/x/term"
)

type Options struct {
	Args struct {
		Intf arke.CANInterfaceName
	} `positional-args:"yes" required:"yes"`
	NoColor bool `long:"no-color"`
}

func execute() error {

	opts := &Options{}

	parser := flags.NewParser(opts, flags.Default)
	_, err := parser.Parse()
	if flags.WroteHelp(err) == true {
		return nil
	}
	if ferr, ok := err.(*flags.Error); ok == true && ferr.Type == flags.ErrRequired {
		parser.WriteHelp(os.Stderr)
		return nil
	}
	if err != nil {
		return err
	}

	if opts.NoColor == true || term.IsTerminal(int(os.Stdout.Fd())) == false {
		for k := range colorCodes {
			colorCodes[k] = ""
		}
	}

	intf, err := socketcan.NewRawInterface(string(opts.Args.Intf))
	if err != nil {
		return err
	}

	frames := make(chan socketcan.CanFrame, 10)
	go func() {
		defer close(frames)
		for {
			f, err := intf.Receive()
			if err != nil {
				if errno, ok := err.(syscall.Errno); ok == true &&
					(errno == syscall.EBADF ||
						errno == syscall.ENETDOWN ||
						errno == syscall.ENODEV) {
					log.Printf("Closed CAN Interface: %s", err)
					return
				}
				log.Printf("Could not receive CAN frame: %s", err)
				continue
			}
			frames <- f
		}

	}()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		os.Exit(0)
	}()

	for f := range frames {
		m, _, err := arke.ParseMessage(&f)
		if err != nil {
			log.Printf("Could not parse CAN Frame: %s", err)
		} else {
			formatMessage(m, f.ID, f.RTR)
		}
	}
	return nil
}

const (
	HEARTBEAT_HEADER int = iota
	HEARTBEAT
	NETWORK_COMMAND_HEADER
	NETWORK_COMMAND
	STANDARD_MESSAGE_HEADER
	STANDARD_MESSAGE
	PRIORITY_MESSAGE_HEADER
	PRIORITY_MESSAGE
	REQUEST_HEADER
	REQUEST
)

var colorCodes = map[int]string{
	HEARTBEAT_HEADER:        "\033[30;45m",
	HEARTBEAT:               "\033[35;49m",
	NETWORK_COMMAND_HEADER:  "\033[30;46m",
	NETWORK_COMMAND:         "\033[36;49m",
	STANDARD_MESSAGE_HEADER: "\033[30;47m",
	STANDARD_MESSAGE:        "\033[m",
	PRIORITY_MESSAGE_HEADER: "\033[39;41m",
	PRIORITY_MESSAGE:        "\033[31;49m",
	REQUEST_HEADER:          "\033[30;43m",
	REQUEST:                 "\033[33;49m",
}

func formatMessage(m arke.ReceivableMessage, idt uint32, RTR bool) {
	now := time.Now().Format(time.RFC3339Nano)
	tpe, cls, id := arke.ExtractCANIDT(idt)
	var header, message int
	switch tpe {
	case arke.StandardMessage:
		header, message = STANDARD_MESSAGE_HEADER, STANDARD_MESSAGE
	case arke.HighPriorityMessage:
		header, message = PRIORITY_MESSAGE_HEADER, PRIORITY_MESSAGE
	default:
	}

	if RTR == true {
		fmt.Printf("%s%s%s %s ID:%d\n", colorCodes[REQUEST_HEADER], now, colorCodes[REQUEST], cls, id)
		return
	}

	if tpe == arke.NetworkControlCommand {
		fmt.Printf("%s%s%s %s\n", colorCodes[NETWORK_COMMAND_HEADER], now, colorCodes[NETWORK_COMMAND], m)
		return
	}

	if tpe == arke.HeartBeat {
		fmt.Printf("%s%s%s %s\n", colorCodes[HEARTBEAT_HEADER], now, colorCodes[HEARTBEAT], m)
		return
	}

	fmt.Printf("%s%s%s ID:%d %s\n", colorCodes[header], now, colorCodes[message], id, m)
}

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}
