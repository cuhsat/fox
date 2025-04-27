package widgets

import (
    "fmt"

    "github.com/cuhsat/fx/internal/fx/args"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/ui/context"
    "github.com/cuhsat/fx/internal/ui/themes"
    "github.com/gdamore/tcell/v2"
)

const (
    filter = "❯"
    follow = "F"
    line = "N"
    wrap = "W"
)

type Prompt struct {
    base
    Lock bool
    Value string
}

func NewPrompt(ctx *context.Context, term tcell.Screen) *Prompt {
    return &Prompt{
        base: base{ctx, term},

        Lock: true,
        Value: "",
    }
}

func (p *Prompt) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    m := p.fmtMode()

    // render blank line
    p.blank(x, y, w, themes.Surface0)

    // render mode
    p.print(x, y, m, themes.Surface3)

    if p.ctx.Mode == mode.Hex {
        return 1
    }

    x += text.Len(m)

    _, heap := hs.Current()

    f := p.fmtFilters(heap)
    s := p.fmtStatus(heap)

    // render filters
    if p.ctx.Mode == mode.Grep || len(f) > 2 {
        p.print(x, y, text.Abr(f, w - (x + text.Len(s))), themes.Surface1)
    }

    // render status
    p.print((w-text.Len(s)), y, s, themes.Surface1)

    if p.Lock {
        p.term.HideCursor()
    } else {
        p.term.ShowCursor(x + text.Len(f)-1, y)
    }

    return 1
}

func (p *Prompt) AddRune(r rune) {
    if !p.Lock {
        p.Value += string(r)
    }
}

func (p *Prompt) DelRune() {
    if !p.Lock && len(p.Value) > 0 {
        p.Value = p.Value[:len(p.Value)-1]
    }
}

func (p *Prompt) Accept() (s string) {
    if !p.Lock {
        s, p.Value = p.Value, ""
    }

    return
}

func (p *Prompt) fmtMode() string {
    return fmt.Sprintf(" %s ", p.ctx.Mode)
}

func (p *Prompt) fmtFilters(h *heap.Heap) (s string) {
    for _, f := range *args.GetFilters() {
        s = fmt.Sprintf("%s %s %s", s, f, filter)
    }

    s = fmt.Sprintf("%s %s ", s, p.Value)

    return 
}

func (p *Prompt) fmtStatus(h *heap.Heap) string {
    f, n, w := "·", "·", "·"

    if p.ctx.Follow {
        f = follow
    }

    if p.ctx.Line {
        n = line
    }

    if p.ctx.Wrap {
        w = wrap
    }

    return fmt.Sprintf(" %d %s%s%s ", len(h.SMap), f, n, w)
}
