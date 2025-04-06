package status

import (
    "github.com/cuhsat/cu/pkg/ui/mode"
)

type Status struct {
    Mode    mode.Mode
    Last    mode.Mode
    
    Numbers bool
    Wrap    bool
}

// singleton
var instance *Status = nil

func NewStatus() *Status {
    if instance == nil {
        instance = &Status{
            Mode: mode.Grep,
            Last: mode.Grep,
            Numbers: true,
            Wrap: false,
        }
    }

    return instance
}

func (s *Status) SwitchMode(m mode.Mode) bool {
    // allow only goto in grep mode
    if m == mode.Goto && s.Mode == mode.Hex {
        return false
    }

    // react only to mode changes
    if m == s.Mode {
        return false
    }

    s.Last = s.Mode
    s.Mode = m

    return true
}

func (s *Status) ToggleNumbers() bool {
    s.Numbers = !s.Numbers

    return s.Numbers
}

func (s *Status) ToggleWrap() bool {
    s.Wrap = !s.Wrap

    return s.Wrap
}
