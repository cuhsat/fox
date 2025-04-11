package widget

import (
    "fmt"
    "strings"

    "github.com/cuhsat/cu/pkg/ui/buffer"
    "github.com/cuhsat/cu/pkg/ui/themes"
)

const (
    Rule = "┃"
)

func (o *Output) hexRender(x, y, w, h int) {
    hb := buffer.Hex(buffer.Context{
        Heap: o.heap,
        Line: o.status.Line,
        Wrap: o.status.Wrap,
        X: o.delta_x,
        Y: o.delta_y,
        W: w,
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
            hex_x += buffer.HexSpace * 2
        }

        // hex values
        o.print(hex_x, line_y, line.Hex, themes.Output)
        text_x := hex_x + len(line.Hex)

        // hex separator
        o.print(text_x, line_y, Rule, themes.Rule)
        text_x += buffer.HexSpace * 2

        // text value
        o.print(text_x, line_y, line.Str, themes.Output)

        // mark found positions
        // for c, f := range heap.Chain {
        //     o.hexMark(hex_x, line_y, c, line.str, f.Name)
        //     o.textMark(text_x, line_y, c, line.str, f.Name)
        // }
    }
}

func (o *Output) hexMark(x, y, c int, s, f string) {
    j := 0

    for j < len(s) {
        i := strings.Index(s[j:], f)

        if i == -1 {
            break
        }

        i += j

        for bx, b := range []byte(f) {
            h := fmt.Sprintf("%02x", b)

            dx := (i*2) + (bx*2)
            dx += dx / 4

            o.print(x + dx, y, h, themes.Colors[c % len(themes.Colors)])
        }

        j = i+1
    }

    return
}
