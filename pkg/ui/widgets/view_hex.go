package widgets

import (
    "github.com/cuhsat/fx/pkg/fx/types/buffer"
    "github.com/cuhsat/fx/pkg/ui/themes"
)

func (v *View) hexRender(x, y, w, h int) {
    rule_w := buffer.SpaceHex * 2

    buf := buffer.Hex(&buffer.Context{
        Heap: v.heap,
        Line: v.ctx.Line,
        Wrap: v.ctx.Wrap,
        X: v.delta_x,
        Y: v.delta_y,
        W: w - (rule_w * 2),
        H: h,
    })

    if len(buf.Lines) > 0 {
        w -= len(buf.Lines[0].Nr) + buffer.SpaceHex
    }

    // set buffer bounds
    v.last_x = max(buf.W, 0)
    v.last_y = max(buf.H - h, 0)

    // render buffer
    for i, line := range buf.Lines {
        line_x := x
        line_y := y + i

        hex_x := line_x

        // offset number
        v.print(hex_x, line_y, line.Nr, themes.Subtext0)
        hex_x += len(line.Nr)

        // offset separator
        v.print(hex_x, line_y, "│", themes.Subtext1)
        hex_x += rule_w

        // hex values
        v.print(hex_x, line_y, line.Hex, themes.Base)
        text_x := hex_x + len(line.Hex)

        // hex separator
        v.print(text_x, line_y, "│", themes.Subtext1)
        text_x += rule_w

        // text value
        v.print(text_x, line_y, line.Str, themes.Base)
    }
}
