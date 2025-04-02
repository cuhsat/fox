package mode

import (
    "strings"
)

type Mode int

const (
    Normal Mode = iota
    Hex
    Shell
)

func (m Mode) String() string {
    modes := [...]string{
        "Normal", 
        "Hex", 
        "Shell", 
    }

    if m < Normal || m > Shell {
      return "..."
    }

    return strings.ToUpper(modes[m])
}
