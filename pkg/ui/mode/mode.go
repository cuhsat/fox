package mode

import (
    "strings"
)

type Mode int

const (
    Normal Mode = iota
    Hex
    Goto
)

func (m Mode) String() string {
    modes := [...]string{
        "Normal", 
        "Hex", 
        "Goto",
    }

    if m < Normal || m > Goto {
      return "..."
    }

    return strings.ToUpper(modes[m])
}
