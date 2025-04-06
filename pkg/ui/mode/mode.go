package mode

import (
    "strings"
)

type Mode int

const (
    Grep Mode = iota
    Hex
    Goto
)

func (m Mode) String() string {
    modes := [...]string{
        "Grep", 
        "Hex", 
        "Goto",
    }

    if m < Grep || m > Goto {
      return "..."
    }

    return strings.ToUpper(modes[m])
}
