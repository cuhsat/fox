package widget

import (
    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/sys/types/block"
)

const (
    rule = "┃"
)

func (o *Output) hexRender(x, y, w, h int) {
    rule_w := block.HexSpace * 2

    hb := block.Hex(block.Context{
        Heap: o.heap,
        Line: o.ctx.Line,
        Wrap: o.ctx.Wrap,
        X: o.delta_x,
        Y: o.delta_y,
        W: w - (rule_w * 2),
        H: h,
    })

    if len(hb.Lines) > 0 {
        w -= len(hb.Lines[0].Nr) + block.HexSpace
    }

    // set block bounds
    o.last_x = max(hb.W, 0)
    o.last_y = max(hb.H - h, 0)

    // render block
    for i, line := range hb.Lines {
        line_x := x
        line_y := y + i

        hex_x := line_x

        if o.ctx.Line {
            // offset number
            o.print(hex_x, line_y, line.Nr, themes.Subtext0)
            hex_x += len(line.Nr)

            // offset separator
            o.print(hex_x, line_y, rule, themes.Subtext1)
            hex_x += rule_w
        }

        // hex values
        o.print(hex_x, line_y, line.Hex, themes.Base)
        text_x := hex_x + len(line.Hex)

        // hex separator
        o.print(text_x, line_y, rule, themes.Subtext1)
        text_x += rule_w

        // text value
        o.print(text_x, line_y, line.Str, themes.Base)
    }
}
