package widget

import (
    "fmt"
    "math"
    "strings"

    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

type Output struct {
    widget

    numbers bool
    max_n   int
    max_x   int
    max_y   int
    delta_x int
    delta_y int
}

func NewOutput(screen tcell.Screen) *Output {
    return &Output{
        widget: widget{
            screen: screen,
        },
        numbers: true,
        max_n: 0,
        max_x: 0,
        max_y: 0,
        delta_x: 0,
        delta_y: 0,
    }
}

func (o *Output) Render(heap *data.Heap, x, y, h int) {
    o.max_x = maxString(heap.SMap)
    o.max_y = len(heap.SMap)
    o.max_n = int(math.Log10(float64(heap.Lines()))) + 1

    for i, se := range heap.SMap[o.delta_y:] {
        if i > h-1 {
            break
        }

        xn := x

        if o.numbers {
            n := fmt.Sprintf("%0*d", o.max_n, se.Nr)
            xn += o.max_n + 1

            o.print(x, y + i, n, theme.Number)
        }

        s := string(heap.MMap[se.Start:se.End])
        d := min(o.delta_x, len(s))

        o.print(xn, y + i, s[d:], theme.Output)

        for z, f := range heap.Chain {
            o.highlight(xn, y + i, z, s[d:], f.Name)
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

func (o *Output) ScrollEnd(page int) {
    o.delta_y = min(o.max_y, o.max_y - page)
}

func (o *Output) ScrollUp(delta int) {
    o.delta_y = max(o.delta_y - delta, 0)
}

func (o *Output) ScrollDown(delta int) {
    o.delta_y = min(o.delta_y + delta, o.max_y-1)
}

func (o *Output) ScrollLeft(delta int) {
    o.delta_x = max(o.delta_x - delta, 0)
}

func (o *Output) ScrollRight(delta int) {
    o.delta_x = min(o.delta_x + delta, o.max_x-1)
}

func (o *Output) ScrollPageUp(page int) {
    o.delta_y = max(o.delta_y - page, 0)
}

func (o *Output) ScrollPageDown(page int) {
    o.delta_y = min(o.delta_y + page, o.max_y)
}

func (o *Output) ToggleNumbers() {
    o.numbers = !o.numbers
}

func (o *Output) highlight(x, y, z int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    o.print(x + i, y, f, theme.Colors[z % len(theme.Colors)])
    
    o.highlight(x + i+1, y, z, s[i+1:], f)
}

func maxString(s data.SMap) (w int) {
    for _, se := range s {
        w = max(w, se.Len)
    }

    return
}
