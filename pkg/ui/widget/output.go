package widget

import (
    "fmt"
    "math"
    "strings"

    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    Space = 1
)

type Output struct {
    widget

    numbers bool
    wrap    bool
    max_x   int
    max_y   int
    last_x  int
    last_y  int
    delta_x int
    delta_y int
}

type dline struct {
    n, s string
}

func NewOutput(screen tcell.Screen) *Output {
    return &Output{
        widget: widget{
            screen: screen,
        },
        numbers: true,
        wrap: false,
        max_x: 0,
        max_y: 0,
        last_x: 0,
        last_y: 0,
        delta_x: 0,
        delta_y: 0,
    }
}

func (o *Output) Render(heap *data.Heap, x, y, w, h int) {
    // set output limits
    o.max_x = wid(heap.SMap)
    o.max_y = len(heap.SMap)

    // convert logical to display lines
    lines := o.buffer(heap, w, h)

    if len(lines) > 0 {
        w -= len(lines[0].n) + Space
    }

    // set buffer bounds
    o.last_x = max(o.max_x - w, 0)
    o.last_y = max(o.max_y - h, 0)

    // render buffer
    for i, line := range lines {
        line_x := x
        line_y := y + i

        // line number
        if o.numbers {
            o.print(line_x, line_y, line.n, theme.Number)
            
            line_x += len(line.n) + Space
        }

        o.print(line_x, line_y, line.s, theme.Output)

        // mark found positions
        for c, f := range heap.Chain {
            o.rmark(line_x, line_y, c, line.s, f.Name)
        }
    }
}

func (o *Output) Reset() {
    o.max_x = 0
    o.max_y = 0
    o.delta_x = 0
    o.delta_y = 0
}

func (o *Output) ScrollBegin() {
    o.delta_y = 0
}

func (o *Output) ScrollEnd() {
    o.delta_y = min(o.max_y, o.last_y)
}

func (o *Output) ScrollUp(delta int) {
    o.delta_y = max(o.delta_y - delta, 0)
}

func (o *Output) ScrollDown(delta int) {
    o.delta_y = min(o.delta_y + delta, o.last_y)
}

func (o *Output) ScrollLeft(delta int) {
    o.delta_x = max(o.delta_x - delta, 0)
}

func (o *Output) ScrollRight(delta int) {
    o.delta_x = min(o.delta_x + delta, o.last_x)
}

func (o *Output) ScrollPageUp(delta int) {
    o.delta_y = max(o.delta_y - delta, 0)
}

func (o *Output) ScrollPageDown(delta int) {
    o.delta_y = min(o.delta_y + delta, o.last_y)
}

func (o *Output) ToggleNumbers() {
    o.numbers = !o.numbers
}

func (o *Output) ToggleWrap() {
    o.wrap = !o.wrap
}

func (o *Output) buffer(heap *data.Heap, w, h int) (l []dline) {
    len_nr := int(math.Log10(float64(heap.Lines()))) + 1

    if o.numbers {
        w -= (len_nr + Space)
    }

    for i, se := range heap.SMap[o.delta_y:] {
        if len(l) >= h {
            return l[:h]
        }

        if i >= h {
            return
        }

        // line number
        n := fmt.Sprintf("%0*d", len_nr, se.Nr)

        // logical line
        s := string(heap.MMap[se.Start:se.End])
        s = s[min(o.delta_x, length(s)):]

        // display lines
        if o.wrap {
            for {
                if length(s) < w+1 {
                    break
                }

                l = append(l, dline{
                    n: n,
                    s: s[:w-1] + "\r",
                })

                s = s[w-1:]
            }
        }

        l = append(l, dline{
            n: n,
            s: s,
        })
    }

    return
}

func (o *Output) rmark(x, y, c int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    o.print(x + i, y, f, theme.Colors[c % len(theme.Colors)])
    
    o.rmark(x + i+1, y, c, s[i+1:], f)
}

func wid(s data.SMap) (w int) {
    for _, se := range s {
        w = max(w, se.Len)
    }

    return
}
