package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	socketcan "github.com/atuleu/golang-socketcan"
	"github.com/formicidae-tracker/libarke/src-go/arke"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Args struct {
		Intf arke.CANInterfaceName
	} `positional-args:"yes" required:"yes"`
}

func execute() error {
	opts := &Options{}
	_, err := flags.Parse(opts)
	if flags.WroteHelp(err) == true {
		return nil
	}
	if err != nil {
		return err
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
		intf.Close()
	}()

	out := log.New(os.Stdout, "", log.LstdFlags)
	for f := range frames {
		m, ID, err := arke.ParseMessage(&f)
		if err != nil {
			log.Printf("Could not parse CAN Frame: %s", err)
		} else {
			out.Printf("ID:%d %s", ID, m.String())
		}
	}
	return nil
}

func main() {
	if err := execute(); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}
