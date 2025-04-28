package context

import (
    "github.com/cuhsat/fx/pkg/fx/types/mode"
    "github.com/cuhsat/fx/pkg/fx/user/config"
    "github.com/gdamore/tcell/v2"
)

type Context struct {
    config.Config

    Root tcell.Screen

    Mode mode.Mode
    Last mode.Mode

    Busy bool
}

func New(root tcell.Screen) *Context {
    return &Context{
        Config: *config.New(),
        Root: root,
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
