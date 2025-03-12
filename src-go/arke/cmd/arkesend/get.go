package main

import "github.com/jessevdk/go-flags"

type GetGroup struct{}

var get = &GetGroup{}

func MustAddCommand(cmd *flags.Command, name, short, long string, data interface{}) *flags.Command {
	res, err := cmd.AddCommand(name, short, long, data)
	if err != nil {
		panic(err.Error())
	}
	return res
}

var getCommand *flags.Command = MustAddCommand(parser.Command, "get", "request data group", "sends a request for some data over the CANbus", get)
