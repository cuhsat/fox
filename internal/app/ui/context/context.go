package context

import (
	"sync"
	"sync/atomic"

	"github.com/gdamore/tcell/v2"

	"github.com/cuhsat/fx/internal/pkg/types/mode"
	"github.com/cuhsat/fx/internal/pkg/user/config"
)

type Context struct {
	sync.RWMutex

	Root tcell.Screen

	cfg *config.Config

	mode mode.Mode
	last mode.Mode

	theme string

	tail atomic.Bool
	line atomic.Bool
	wrap atomic.Bool
}

func New(root tcell.Screen) *Context {
	cfg := config.New()
	ctx := &Context{
		// screen
		Root: root,

		// config
		cfg: cfg,

		// modes
		mode: mode.Default,
		last: mode.Default,

		// theme
		theme: cfg.Theme,
	}

	// flags
	ctx.tail.Store(cfg.Tail)
	ctx.line.Store(cfg.Line)
	ctx.wrap.Store(cfg.Wrap)

	return ctx
}

func (ctx *Context) Mode() mode.Mode {
	ctx.RLock()
	defer ctx.RUnlock()
	return ctx.mode
}

func (ctx *Context) Last() mode.Mode {
	ctx.RLock()
	defer ctx.RUnlock()
	return ctx.last
}

func (ctx *Context) Theme() string {
	ctx.RLock()
	defer ctx.RUnlock()
	return ctx.theme
}

func (ctx *Context) IsTail() bool {
	return ctx.tail.Load()
}

func (ctx *Context) IsLine() bool {
	return ctx.line.Load()
}

func (ctx *Context) IsWrap() bool {
	return ctx.wrap.Load()
}

func (ctx *Context) Interrupt() {
	ctx.Root.PostEvent(tcell.NewEventInterrupt(nil))
}

func (ctx *Context) SwitchMode(m mode.Mode) bool {
	// deny goto in hex mode
	if m == mode.Goto && ctx.Mode() == mode.Hex {
		return false
	}

	// react only to mode changes
	if m == ctx.Mode() {
		return false
	}

	ctx.Lock()
	ctx.last = ctx.mode
	ctx.mode = m
	ctx.Unlock()

	return true
}

func (ctx *Context) ChangeTheme(t string) {
	ctx.Lock()
	ctx.theme = t
	ctx.Unlock()
}

func (ctx *Context) ToggleFollow() {
	ctx.tail.Store(!ctx.tail.Load())
}

func (ctx *Context) ToggleNumbers() {
	ctx.line.Store(!ctx.line.Load())
}

func (ctx *Context) ToggleWrap() {
	ctx.wrap.Store(!ctx.wrap.Load())
}

func (ctx *Context) Background(fn func()) {
	go func() {
		fn()
		ctx.Interrupt()
	}()
}

func (ctx *Context) Save() {
	ctx.cfg.Theme = ctx.Theme()
	ctx.cfg.Tail = ctx.IsTail()
	ctx.cfg.Line = ctx.IsLine()
	ctx.cfg.Wrap = ctx.IsWrap()
	ctx.cfg.Save()
}
