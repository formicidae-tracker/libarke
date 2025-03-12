package arke

import (
	"strings"

	socketcan "github.com/atuleu/golang-socketcan"
	"github.com/jessevdk/go-flags"
)

type CANInterfaceName string

func (n *CANInterfaceName) Complete(match string) []flags.Completion {
	availables, err := socketcan.ListCANInterfaces()
	completions := make([]flags.Completion, 0, len(availables))
	if err != nil {
		return nil
	}

	match = strings.ToLower(match)
	for _, name := range availables {
		if strings.HasPrefix(name, match) {
			completions = append(completions, flags.Completion{Item: name})
		}
	}

	return completions
}
