package lib

import (
    "strconv"
    "strings"

    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/block"
    "github.com/cuhsat/fx/internal/ui/themes"
)

func (b *Buffer) textRender(x, y, w, h int) {
    tb := block.Text(block.Context{
        Heap: b.heap,
        Line: b.ctx.Line,
        Wrap: b.ctx.Wrap,
        X: b.delta_x,
        Y: b.delta_y,
        W: w,
        H: h,
    })

    b.smap = tb.SMap

    if len(tb.Lines) > 0 {
        w -= len(tb.Lines[0].Nr) + block.TextSpace
    }

    // set block bounds
    b.last_x = max(tb.W - w, 0)
    b.last_y = max(tb.H - h, 0)

    // render block
    for i, line := range tb.Lines {
        line_x := x
        line_y := y + i

        // line number
        if b.ctx.Line {
            b.print(line_x, line_y, line.Nr, themes.Subtext0)
            line_x += len(line.Nr) + block.TextSpace
        }

        // text value
        b.print(line_x, line_y, line.Str, themes.Base)

        // mark found positions
        for c, f := range *types.GetFilters() {
            b.textMark(line_x, line_y, c, line.Str, f)
        }
    }
}

func (b *Buffer) textMark(x, y, c int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    len_i := text.Len(s[:i])

    b.print(x + len_i, y, f, themes.Colors[c % len(themes.Colors)])
    
    b.textMark(x + len_i+1, y, c, s[i+1:], f)
}

func (b *Buffer) textGoto(s string) {
    var nr int

    switch s[0] {
    case '+':
        delta, _ := strconv.Atoi(s[1:])
        nr = b.smap[b.delta_y].Nr + delta

    case '-':
        delta, _ := strconv.Atoi(s[1:])
        nr = b.smap[b.delta_y].Nr - delta

    default:
        nr, _ = strconv.Atoi(s)
    }

    y := b.smap.Find(nr)

    if y >= 0 {
        b.ScrollTo(b.delta_x, y)
    }
}
