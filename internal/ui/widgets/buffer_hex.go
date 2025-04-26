package widgets

import (
    "github.com/cuhsat/fx/internal/fx/types/layer"
    "github.com/cuhsat/fx/internal/ui/themes"
)

func (b *Buffer) hexRender(x, y, w, h int) {
    rule_w := layer.HexSpace * 2

    hl := layer.Hex(&layer.Context{
        Heap: b.heap,
        Line: b.ctx.Line,
        Wrap: b.ctx.Wrap,
        X: b.delta_x,
        Y: b.delta_y,
        W: w - (rule_w * 2),
        H: h,
    })

    if len(hl.Lines) > 0 {
        w -= len(hl.Lines[0].Nr) + layer.HexSpace
    }

    // set layer bounds
    b.last_x = max(hl.W, 0)
    b.last_y = max(hl.H - h, 0)

    // render layer
    for i, line := range hl.Lines {
        line_x := x
        line_y := y + i

        hex_x := line_x

        // offset number
        b.print(hex_x, line_y, line.Nr, themes.Subtext0)
        hex_x += len(line.Nr)

        // offset separator
        b.print(hex_x, line_y, "│", themes.Subtext1)
        hex_x += rule_w

        // hex values
        b.print(hex_x, line_y, line.Hex, themes.Base)
        text_x := hex_x + len(line.Hex)

        // hex separator
        b.print(text_x, line_y, "│", themes.Subtext1)
        text_x += rule_w

        // text value
        b.print(text_x, line_y, line.Str, themes.Base)
    }
}
