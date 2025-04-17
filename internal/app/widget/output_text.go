package widget

import (
    "strconv"
    "strings"

    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/cuhsat/fx/internal/sys/types"
    "github.com/cuhsat/fx/internal/sys/types/block"
)

func (o *Output) textRender(x, y, w, h int) {
    tb := block.Text(block.Context{
        Heap: o.heap,
        Line: o.ctx.Line,
        Wrap: o.ctx.Wrap,
        X: o.delta_x,
        Y: o.delta_y,
        W: w,
        H: h,
    })

    o.smap = tb.SMap

    if len(tb.Lines) > 0 {
        w -= len(tb.Lines[0].Nr) + block.TextSpace
    }

    // set block bounds
    o.last_x = max(tb.W - w, 0)
    o.last_y = max(tb.H - h, 0)

    // render block
    for i, line := range tb.Lines {
        line_x := x
        line_y := y + i

        // line number
        if o.ctx.Line {
            o.print(line_x, line_y, line.Nr, themes.Subtext0)
            line_x += len(line.Nr) + block.TextSpace
        }

        // text value
        o.print(line_x, line_y, line.Str, themes.Base)

        // mark found positions
        for c, f := range *types.GetFilters() {
            o.textMark(line_x, line_y, c, line.Str, f)
        }
    }
}

func (o *Output) textMark(x, y, c int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    len_i := text.Len(s[:i])

    o.print(x + len_i, y, f, themes.Colors[c % len(themes.Colors)])
    
    o.textMark(x + len_i+1, y, c, s[i+1:], f)
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
