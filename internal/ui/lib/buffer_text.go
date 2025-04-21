package lib

import (
    "strconv"

    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types/layers"
    "github.com/cuhsat/fx/internal/ui/themes"
)

func (b *Buffer) textRender(x, y, w, h int) {
    tl := layers.Text(&layers.Context{
        Heap: b.heap,
        Line: b.ctx.Line,
        Wrap: b.ctx.Wrap,
        X: b.delta_x,
        Y: b.delta_y,
        W: w,
        H: h,
    })

    b.smap = tl.SMap

    nr_w := text.Dec(b.heap.Length()) + layers.TextSpace

    if b.ctx.Line {
        w -= nr_w
    }

    // set layer bounds
    b.last_x = max(tl.W - w, 0)
    b.last_y = max(tl.H - h, 0)

    // render lines
    for i, line := range tl.Lines {
        line_x := x
        line_y := y + i

        // line number
        if b.ctx.Line {
            b.print(line_x, line_y, line.Nr, themes.Subtext0)
            line_x += nr_w
        }

        // text value
        if len(line.Str) > 0 {
            b.print(line_x, line_y, line.Str, themes.Base)
        }
    }

    // render parts on top
    for _, part := range tl.Parts {
        part_x := x + part.X
        part_y := y + part.Y

        if b.ctx.Line {
            part_x += nr_w
        }

        // part value
        b.print(part_x, part_y, part.Str, themes.Subtext2)
    }
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
