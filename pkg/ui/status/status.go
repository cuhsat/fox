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

        // defaults
        instance = &Status{
            Mode: mode.Less,
            Last: mode.Less,
            Numbers: true,
            Wrap: false,
        }
    }

    return instance
}

func (s *Status) SwitchMode(m mode.Mode) bool {
    // deny goto in hex mode
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

func (s *Status) ToggleNumbers() {
    s.Numbers = !s.Numbers
}

func (s *Status) ToggleWrap() {
    s.Wrap = !s.Wrap
}
