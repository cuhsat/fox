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
)

type Mode int

func (m Mode) String() string {
    modes := [...]string{
        "Less",
        "Grep", 
        "Hex", 
        "Goto",
        "Open",
    }

    if int(m) < 0 || int(m) > len(modes) {
        return "..."
    }

    return strings.ToUpper(modes[m])
}
