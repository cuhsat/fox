package mode

import (
	"strings"
)

const (
	Default = Less
)

const (
	Less Mode = iota
	Grep
	Hex
	Goto
	Open
	User
)

type Mode int

func (m Mode) String() string {
	modes := [...]string{
		"Less",
		"Grep",
		"Hex",
		"Goto",
		"Open",
		"User",
	}

	if int(m) < 0 || int(m) > len(modes) {
		return "..."
	}

	return strings.ToUpper(modes[m])
}

func (m Mode) Interactive() bool {
	switch m {
	case Grep, Goto, Open, User:
		return true
	default:
		return false
	}
}
