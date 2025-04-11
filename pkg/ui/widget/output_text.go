package widget

import (
    "strconv"
    "strings"

    "github.com/cuhsat/cu/pkg/ui/buffer"
    "github.com/cuhsat/cu/pkg/ui/themes"
)

func (o *Output) textRender(x, y, w, h int) {
    tb := buffer.Text(buffer.Context{
        Heap: o.heap,
        Line: o.status.Line,
        Wrap: o.status.Wrap,
        X: o.delta_x,
        Y: o.delta_y,
        W: w,
        H: h,
    })

    o.smap = tb.SMap

    if len(tb.Lines) > 0 {
        w -= len(tb.Lines[0].Nr) + buffer.TextSpace
    }

    // set buffer bounds
    o.last_x = max(tb.W - w, 0)
    o.last_y = max(tb.H - h, 0)

    // render buffer
    for i, line := range tb.Lines {
        line_x := x
        line_y := y + i

        // line number
        if o.status.Line {
            o.print(line_x, line_y, line.Nr, themes.Hint)
            line_x += len(line.Nr) + buffer.TextSpace
        }

        // text value
        o.print(line_x, line_y, line.Str, themes.Output)

        // mark found positions
        for c, f := range o.heap.Chain {
            o.textMark(line_x, line_y, c, line.Str, f.Name)
        }
    }
}

func (o *Output) textMark(x, y, c int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    o.print(x + i, y, f, themes.Colors[c % len(themes.Colors)])
    
    o.textMark(x + i+1, y, c, s[i+1:], f)
}

func (o *Output) textGoto(s string) {
    var nr int

    switch s[0] {
    case '+':
        delta, _ := strconv.Atoi(s[1:])
        nr = o.smap[o.delta_y].Nr + delta

    case '-':
        delta, _ := strconv.Atoi(s[1:])
        nr = o.smap[o.delta_y].Nr - delta

    default:
        nr, _ = strconv.Atoi(s)
    }

    y := o.smap.Find(nr)

    if y >= 0 {
        o.ScrollTo(o.delta_x, y)
    }
}
