package widget

import (
    "fmt"

    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/sys/heap"
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/cuhsat/fx/internal/sys/types"
    "github.com/cuhsat/fx/internal/sys/types/mode"
    "github.com/gdamore/tcell/v2"
)

const (
    Separator = "❯"
    Cursor = "_"
    Follow = "F"
    Line = "N"
    Wrap = "W"
)

type Prompt struct {
    widget

    Lock bool
    Value string
}

func NewPrompt(screen tcell.Screen, status *Status) *Prompt {
    return &Prompt{
        widget: widget{
            screen: screen,
            status: status,
        },

        Lock: true,
        Value: "",
    }
}

func (p *Prompt) Render(hs *heapset.HeapSet, x, y, w, h int) int {
    m := p.formatMode()

    // render blank line
    p.printBlank(x, y, w, themes.Line)

    // render mode
    p.print(x, y, m, themes.Mode)

    if p.status.Mode == mode.Hex {
        return 1
    }

    x += text.Length(m)

    _, heap := hs.Current()

    f := p.formatFilters(heap)
    s := p.formatStatus(heap)

    // render filters
    if !p.Lock {
        p.print(x, y, text.Abbrev(f, x, w-text.Length(s)), themes.Input)
    }

    // render status
    p.print(w-text.Length(s), y, s, themes.Input)

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
    return fmt.Sprintf(" %s ", p.status.Mode)
}

func (p *Prompt) formatFilters(h *heap.Heap) (s string) {
    if p.status.Mode == mode.Grep {
        for _, f := range *types.GetFilters() {
            s = fmt.Sprintf("%s %s %s", s, f, Separator)
        }        
    }

    s = fmt.Sprintf("%s %s%s ", s, p.Value, Cursor)

    return 
}

func (p *Prompt) formatStatus(h *heap.Heap) string {
    f, n, w := " ", " ", " "

    if p.status.Follow {
        f = Follow
    }

    if p.status.Line {
        n = Line
    }

    if p.status.Wrap {
        w = Wrap
    }

    return fmt.Sprintf(" %d ∣ %s ∣ %s ∣ %s ", len(h.SMap), f, n, w)
}
