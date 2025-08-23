package context

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gdamore/tcell/v2"

	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
	"github.com/cuhsat/fox/internal/pkg/user/config"
)

type Context struct {
	sync.RWMutex

	Root tcell.Screen

	mode mode.Mode
	last mode.Mode

	model string
	theme string

	n atomic.Bool
	w atomic.Bool
	t atomic.Bool
	p atomic.Bool
}

func New(root tcell.Screen) *Context {
	cfg := config.Get()
	ctx := &Context{
		// screen
		Root: root,

		// modes
		mode: mode.Default,
		last: mode.Default,

		// model
		model: cfg.GetString("ai.model"),

		// theme
		theme: cfg.GetString("ui.theme"),
	}

	ctx.n.Store(cfg.GetBool("ui.state.n"))
	ctx.w.Store(cfg.GetBool("ui.state.w"))
	ctx.t.Store(cfg.GetBool("ui.state.t"))
	ctx.p.Store(false)

	s := strings.ToUpper(flags.Get().UI.State)

	// precede flags
	if strings.ContainsRune(s, '-') {
		ctx.n.Store(false)
		ctx.w.Store(false)
		ctx.t.Store(false)
	} else if len(s) > 0 {
		ctx.n.Store(strings.ContainsRune(s, 'N'))
		ctx.w.Store(strings.ContainsRune(s, 'W'))
		ctx.t.Store(strings.ContainsRune(s, 'T'))
	}

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

func (ctx *Context) IsNavi() bool {
	return ctx.n.Load()
}

func (ctx *Context) IsWrap() bool {
	return ctx.w.Load()
}

func (ctx *Context) IsFollow() bool {
	return ctx.t.Load()
}

func (ctx *Context) IsPinned() bool {
	return ctx.p.Load()
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

func (ctx *Context) ToggleNavi() {
	ctx.n.Store(!ctx.n.Load())
}

func (ctx *Context) ToggleWrap() {
	ctx.w.Store(!ctx.w.Load())
}

func (ctx *Context) ToggleFollow() {
	ctx.t.Store(!ctx.t.Load())
}

func (ctx *Context) TogglePinned() {
	ctx.p.Store(!ctx.p.Load())
}

func (ctx *Context) Background(fn func()) {
	go func() {
		fn()
		ctx.ForceRender()
	}()
}

func (ctx *Context) Save() {
	cfg := config.Get()

	cfg.Set("ai.model", ctx.Model())
	cfg.Set("ui.theme", ctx.Theme())
	cfg.Set("ui.state.n", ctx.IsNavi())
	cfg.Set("ui.state.w", ctx.IsWrap())
	cfg.Set("ui.state.t", ctx.IsFollow())

	if !flags.Get().Opt.Readonly {
		config.Save()
	}
}
