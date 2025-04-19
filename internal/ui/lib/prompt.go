package lib

import (
    "fmt"

    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/ui/themes"
    "github.com/gdamore/tcell/v2"
)

const (
    cursor = "_"
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

func NewPrompt(ctx *Context, term tcell.Screen) *Prompt {
    return &Prompt{
        base: base{
            ctx: ctx,
            term: term,
        },

        Lock: true,
        Value: "",
    }
}

func (p *Prompt) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    m := p.formatMode()

    // render blank line
    p.blank(x, y, w, themes.Surface0)

    // render mode
    p.print(x, y, m, themes.Surface3)

    if p.ctx.Mode == mode.Hex {
        return 1
    }

    x += text.Len(m)

    _, heap := hs.Current()

    f := p.formatFilters(heap)
    s := p.formatStatus(heap)

    // render filters
    if !p.Lock {
        p.print(x, y, text.Abr(f, w - (x + text.Len(s))), themes.Surface1)
    }

    // render status
    p.print(w-text.Len(s), y, s, themes.Surface1)

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

func (p *Prompt) formatMode() string {
    return fmt.Sprintf(" %s ", p.ctx.Mode)
}

func (p *Prompt) formatFilters(h *heap.Heap) (s string) {
    if p.ctx.Mode == mode.Grep {
        for _, f := range *types.GetFilters() {
            s = fmt.Sprintf("%s %s %s", s, f, filter)
        }        
    }

    s = fmt.Sprintf("%s %s%s ", s, p.Value, cursor)

    return 
}

func (p *Prompt) formatStatus(h *heap.Heap) string {
    f, n, w := " ", " ", " "

    if p.ctx.Follow {
        f = follow
    }

    if p.ctx.Line {
        n = line
    }

    if p.ctx.Wrap {
        w = wrap
    }

    return fmt.Sprintf(" %d ∣ %s ∣ %s ∣ %s ", len(h.SMap), f, n, w)
}
