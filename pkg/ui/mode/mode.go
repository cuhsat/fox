package mode

import (
    "strings"
)

type Mode int

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

    if m < Less || m > Goto {
      return "..."
    }

    return strings.ToUpper(modes[m])
}
