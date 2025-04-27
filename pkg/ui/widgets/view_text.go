package widgets

import (
    "strconv"

    "github.com/cuhsat/fx/pkg/fx/text"
    "github.com/cuhsat/fx/pkg/fx/types/buffer"
    "github.com/cuhsat/fx/pkg/ui/themes"
)

func (v *View) textRender(x, y, w, h int) {
    buf := buffer.Text(&buffer.Context{
        Heap: v.heap,
        Line: v.ctx.Line,
        Wrap: v.ctx.Wrap,
        X: v.delta_x,
        Y: v.delta_y,
        W: w,
        H: h,
    })

    v.smap = buf.SMap

    if v.ctx.Line {
        w -= text.Dec(v.heap.Length()) + buffer.SpaceText
    }

    // set buffer bounds
    v.last_x = max(buf.W - w, 0)
    v.last_y = max(buf.H - h, 0)

    // render lines
    for i, line := range buf.Lines {
        line_x := x
        line_y := y + i

        // line number
        if v.ctx.Line {
            v.print(line_x, line_y, line.Nr, themes.Subtext0)
            line_x += len(line.Nr) + buffer.SpaceText
        }

        // text value
        if len(line.Str) > 0 {
            v.print(line_x, line_y, line.Str, themes.Base)
        }
    }

    // render parts on top
    for _, part := range buf.Parts {
        part_x := x + part.X
        part_y := y + part.Y

        if v.ctx.Line {
            part_x += len(buf.Lines[0].Nr) + buffer.SpaceText
        }

        // part value
        v.print(part_x, part_y, part.Str, themes.Subtext2)
    }
}

func (v *View) textGoto(s string) {
    var nr int

    switch s[0] {
    case '+':
        delta, _ := strconv.Atoi(s[1:])
        nr = v.smap[v.delta_y].Nr + delta

    case '-':
        delta, _ := strconv.Atoi(s[1:])
        nr = v.smap[v.delta_y].Nr - delta

    default:
        nr, _ = strconv.Atoi(s)
    }

    y := v.smap.Find(nr)

    if y >= 0 {
        v.ScrollTo(v.delta_x, y)
    }
}
