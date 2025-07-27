package context

import (
	"sync"
	"sync/atomic"

	"github.com/gdamore/tcell/v2"

	"github.com/hiforensics/fox/internal/pkg/types/mode"
	"github.com/hiforensics/fox/internal/pkg/user/config"
)

type Context struct {
	sync.RWMutex

	Root tcell.Screen

	cfg *config.Config

	mode mode.Mode
	last mode.Mode

	model string
	theme string

	follow  atomic.Bool
	numbers atomic.Bool
	wrap    atomic.Bool
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

		// model
		model: cfg.Model,

		// theme
		theme: cfg.Theme,
	}

	// flags
	ctx.follow.Store(cfg.Follow)
	ctx.numbers.Store(cfg.Numbers)
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

func (ctx *Context) Model() string {
	ctx.RLock()
	defer ctx.RUnlock()
	return ctx.model
}

func (ctx *Context) Theme() string {
	ctx.RLock()
	defer ctx.RUnlock()
	return ctx.theme
}

func (ctx *Context) IsFollow() bool {
	return ctx.follow.Load()
}

func (ctx *Context) IsNumbers() bool {
	return ctx.numbers.Load()
}

func (ctx *Context) IsWrap() bool {
	return ctx.wrap.Load()
}

func (ctx *Context) ForceRender() {
	_ = ctx.Root.PostEvent(tcell.NewEventInterrupt(nil))
}

func (ctx *Context) SwitchMode(m mode.Mode) bool {
	// deny goto in static modes
	if m == mode.Goto && ctx.Mode().Static() {
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
	ctx.follow.Store(!ctx.follow.Load())
}

func (ctx *Context) ToggleNumbers() {
	ctx.numbers.Store(!ctx.numbers.Load())
}

func (ctx *Context) ToggleWrap() {
	ctx.wrap.Store(!ctx.wrap.Load())
}

func (ctx *Context) Background(fn func()) {
	go func() {
		fn()
		ctx.ForceRender()
	}()
}

func (ctx *Context) Save() {
	ctx.cfg.Model = ctx.Model()
	ctx.cfg.Theme = ctx.Theme()
	ctx.cfg.Follow = ctx.IsFollow()
	ctx.cfg.Numbers = ctx.IsNumbers()
	ctx.cfg.Wrap = ctx.IsWrap()
	ctx.cfg.Save()
}
