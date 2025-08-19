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

	cfg *config.Config

	mode mode.Mode
	last mode.Mode

	model string
	theme string

	t atomic.Bool
	n atomic.Bool
	w atomic.Bool
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

	ctx.t.Store(cfg.Follow)
	ctx.n.Store(cfg.Numbers)
	ctx.w.Store(cfg.Wrap)

	ctx.Precede()

	return ctx
}

func (ctx *Context) Precede() {
	flg := flags.Get()

	s := strings.ToUpper(flg.UI.State)

	// overwrite flags
	if strings.ContainsRune(s, '-') {
		ctx.t.Store(false)
		ctx.n.Store(false)
		ctx.w.Store(false)
	} else if len(flg.UI.State) > 0 {
		ctx.t.Store(strings.ContainsRune(s, 'T'))
		ctx.n.Store(strings.ContainsRune(s, 'N'))
		ctx.w.Store(strings.ContainsRune(s, 'W'))
	}

	// overwrite theme
	if len(flg.UI.Theme) > 0 {
		ctx.theme = flg.UI.Theme
	}

	// overwrite model
	if len(flg.AI.Model) > 0 {
		ctx.model = flg.AI.Model
	}
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
	return ctx.t.Load()
}

func (ctx *Context) IsNumbers() bool {
	return ctx.n.Load()
}

func (ctx *Context) IsWrap() bool {
	return ctx.w.Load()
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
	ctx.t.Store(!ctx.t.Load())
}

func (ctx *Context) ToggleNumbers() {
	ctx.n.Store(!ctx.n.Load())
}

func (ctx *Context) ToggleWrap() {
	ctx.w.Store(!ctx.w.Load())
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

	if !flags.Get().Opt.Readonly {
		ctx.cfg.Save()
	}
}
