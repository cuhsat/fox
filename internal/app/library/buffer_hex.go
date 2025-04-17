package library

import (
    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/sys/types/block"
)

const (
    rule = "┃"
)

func (b *Buffer) hexRender(x, y, w, h int) {
    rule_w := block.HexSpace * 2

    hb := block.Hex(block.Context{
        Heap: b.heap,
        Line: b.ctx.Line,
        Wrap: b.ctx.Wrap,
        X: b.delta_x,
        Y: b.delta_y,
        W: w - (rule_w * 2),
        H: h,
    })

    if len(hb.Lines) > 0 {
        w -= len(hb.Lines[0].Nr) + block.HexSpace
    }

    // set block bounds
    b.last_x = max(hb.W, 0)
    b.last_y = max(hb.H - h, 0)

    // render block
    for i, line := range hb.Lines {
        line_x := x
        line_y := y + i

        hex_x := line_x

        if b.ctx.Line {
            // offset number
            b.print(hex_x, line_y, line.Nr, themes.Subtext0)
            hex_x += len(line.Nr)

            // offset separator
            b.print(hex_x, line_y, rule, themes.Subtext1)
            hex_x += rule_w
        }

        // hex values
        b.print(hex_x, line_y, line.Hex, themes.Base)
        text_x := hex_x + len(line.Hex)

        // hex separator
        b.print(text_x, line_y, rule, themes.Subtext1)
        text_x += rule_w

        // text value
        b.print(text_x, line_y, line.Str, themes.Base)
    }
}
