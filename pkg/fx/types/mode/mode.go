package mode

import (
	"strings"
)

const (
	Default = Less
)

const (
	Less = Mode("Less")
	Grep = Mode("Grep")
	Hex  = Mode("Hex")
	Goto = Mode("Goto")
	Open = Mode("Open")
)

type Mode string

func (m Mode) String() string {
	return strings.ToUpper(string(m))
}

func (m Mode) Prompt() bool {
	switch m {
	case Less, Hex:
		return false
	default:
		return true
	}
}
