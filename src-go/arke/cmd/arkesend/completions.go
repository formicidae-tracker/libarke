package main

import (
	"strings"

	"github.com/formicidae-tracker/libarke/src-go/arke"
	"github.com/jessevdk/go-flags"
)

type NodeClassName string

var nodeClassName = map[string]arke.NodeClass{
	"broadcast": arke.BroadcastClass,
	"zeus":      arke.ZeusClass,
	"helios":    arke.HeliosClass,
	"celaeno":   arke.CelaenoClass,
	"notus":     arke.NotusClass,
}

func (c *NodeClassName) Complete(match string) []flags.Completion {
	match = strings.ToLower(match)
	completions := make([]flags.Completion, 0, len(nodeClassName))
	for name, _ := range nodeClassName {
		if strings.HasPrefix(name, match) == true {
			completions = append(completions, flags.Completion{
				Item: name,
			})
		}
	}
	return completions
}

func (c *NodeClassName) Class() arke.NodeClass {
	if c, ok := nodeClassName[string(*c)]; ok == true {
		return c
	}
	return arke.NodeClass(arke.NodeClassMask)
}
