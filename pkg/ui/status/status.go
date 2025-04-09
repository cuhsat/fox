package status

import (
    "github.com/cuhsat/cu/pkg/fs/config"
    "github.com/cuhsat/cu/pkg/ui/mode"
)

type Status struct {
    Mode mode.Mode
    Last mode.Mode
    
    Follow bool

    Line bool
    Wrap bool
}

func NewStatus(f bool) *Status {
    cfg := config.GetConfig()

    return &Status{
        Mode: mode.Less,
        Last: mode.Less,

        // init from flags
        Follow: f,

        // init from config
        Line: cfg.UI.Line,
        Wrap: cfg.UI.Wrap,
    }
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

func (s *Status) ToggleFollow() {
    s.Follow = !s.Follow
}

func (s *Status) ToggleNumbers() {
    s.Line = !s.Line
}

func (s *Status) ToggleWrap() {
    s.Wrap = !s.Wrap
}
