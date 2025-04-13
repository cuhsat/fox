package widget

import (
    "github.com/cuhsat/cu/internal/sys/files/config"
    "github.com/cuhsat/cu/internal/sys/types/mode"
)

type Status struct {
    config.Config

    Mode mode.Mode
    Last mode.Mode
}

func NewStatus() *Status {
    return &Status{
        Config: *config.Load(),

        Mode: mode.Default,
        Last: mode.Default,
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
