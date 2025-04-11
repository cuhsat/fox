package status

import (
    "github.com/cuhsat/cu/pkg/fs/config"
    "github.com/cuhsat/cu/pkg/ui/mode"
)

const (
    DefaultMode = mode.Less
)

type Status struct {
    Mode mode.Mode
    Last mode.Mode

    config.Config
}

func NewStatus() *Status {
    return &Status{
        Mode: DefaultMode,
        Last: DefaultMode,

        Config: *config.GetConfig(),
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
