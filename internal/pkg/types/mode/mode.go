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
	Rag  = Mode("RAG")
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
