package arke

import (
	"fmt"
	"strings"
)

var nameByClass = make(map[NodeClass]string)

var classByName = make(map[string]NodeClass)

func ClassName(c NodeClass) string {
	if n, ok := nameByClass[c]; ok == true {
		return n
	}
	return "<unknown>"
}

func Class(s string) (NodeClass, error) {
	if c, ok := classByName[strings.ToLower(s)]; ok == true {
		return c, nil
	}
	return 0, fmt.Errorf("Unknown node class '%s'", s)
}

func init() {
	nameByClass[ZeusClass] = "Zeus"
	nameByClass[CelaenoClass] = "Celaeno"
	nameByClass[HeliosClass] = "Helios"

	for c, n := range nameByClass {
		classByName[strings.ToLower(n)] = c
	}
}
