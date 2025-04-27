package context

import (
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/fx/user/config"
)

type Context struct {
    config.Config

    Mode mode.Mode
    Last mode.Mode

    Busy bool
}

func New() *Context {
    return &Context{
        Config: *config.New(),
        Mode: mode.Default,
        Last: mode.Default,
        Busy: false,
    }
}

func (ctx *Context) SwitchMode(m mode.Mode) bool {
    // deny goto in hex mode
    if m == mode.Goto && ctx.Mode == mode.Hex {
        return false
    }

    // react only to mode changes
    if m == ctx.Mode {
        return false
    }

    ctx.Last = ctx.Mode
    ctx.Mode = m

    return true
}

func (ctx *Context) ToggleFollow() {
    ctx.Follow = !ctx.Follow
}

func (ctx *Context) ToggleNumbers() {
    ctx.Line = !ctx.Line
}

func (ctx *Context) ToggleWrap() {
    ctx.Wrap = !ctx.Wrap
}
