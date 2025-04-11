package widget

import (
    "github.com/cuhsat/cu/pkg/ui/buffer"
    "github.com/cuhsat/cu/pkg/ui/themes"
)

const (
    Rule = "┃"
)

func (o *Output) hexRender(x, y, w, h int) {
    rule_w := buffer.HexSpace * 2

    hb := buffer.Hex(buffer.Context{
        Heap: o.heap,
        Line: o.status.Line,
        Wrap: o.status.Wrap,
        X: o.delta_x,
        Y: o.delta_y,
        W: w - (rule_w * 2),
        H: h,
    })

    if len(hb.Lines) > 0 {
        w -= len(hb.Lines[0].Nr) + buffer.HexSpace
    }

    // set buffer bounds
    o.last_x = max(hb.W, 0)
    o.last_y = max(hb.H - h, 0)

    // render buffer
    for i, line := range hb.Lines {
        line_x := x
        line_y := y + i

        hex_x := line_x

        if o.status.Line {
            // offset number
            o.print(hex_x, line_y, line.Nr, themes.Hint)
            hex_x += len(line.Nr)

            // offset separator
            o.print(hex_x, line_y, Rule, themes.Rule)
            hex_x += rule_w
        }

        // hex values
        o.print(hex_x, line_y, line.Hex, themes.Output)
        text_x := hex_x + len(line.Hex)

        // hex separator
        o.print(text_x, line_y, Rule, themes.Rule)
        text_x += rule_w

        // text value
        o.print(text_x, line_y, line.Str, themes.Output)
    }
}
