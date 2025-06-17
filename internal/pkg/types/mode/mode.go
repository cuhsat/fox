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
	Goto = Mode("Goto")
	Open = Mode("Open")
	Hex  = Mode("Hex")
)

type Mode string

func (m Mode) String() string {
	return strings.ToUpper(string(m))
}

func (m Mode) Filter() bool {
	switch m {
	case Less, Grep:
		return true
	default:
		return false
	}
}

func (m Mode) Prompt() bool {
	switch m {
	case Less, Hex:
		return false
	default:
		return true
	}
}

func (m Mode) Static() bool {
	switch m {
	case Hex:
		return true
	default:
		return false
	}
}
