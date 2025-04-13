package mode

import (
    "strings"
)

type Mode int

const (
    Default = Less
)

const (
    Less Mode = iota
    Grep
    Hex
    Goto
)

func (m Mode) String() string {
    modes := [...]string{
        "Less",
        "Grep", 
        "Hex", 
        "Goto",
    }

    if int(m) < 0 || int(m) > len(modes) {
      return "..."
    }

    return strings.ToUpper(modes[m])
}
