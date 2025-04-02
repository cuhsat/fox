package widget

import (
    "fmt"
    "strings"

    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/ui/theme"
)

const (
    HexSpace = 1
)

const (
    Rule = "┃"
)

type hexLine struct {
    off, hex, str string
}

func (o *Output) hexRender(heap *heap.Heap, x, y, w, h int) {
    // convert logical to display lines
    lines, max_y := o.hexBuffer(heap, w, h)

    if len(lines) > 0 {
        w -= len(lines[0].off) + HexSpace
    }

    // set buffer bounds
    o.last_x = max(w, 0)
    o.last_y = max(max_y - h, 0)

    // render buffer
    for i, line := range lines {
        line_x := x
        line_y := y + i

        // offset number
        o.print(line_x, line_y, line.off, theme.Hint)
        hex_x := line_x + len(line.off)

        // offset separator
        o.print(hex_x, line_y, Rule, theme.Rule)
        hex_x += HexSpace * 2

        // hex values
        o.print(hex_x, line_y, line.hex, theme.Output)
        text_x := hex_x + len(line.hex)

        // hex separator
        o.print(text_x, line_y, Rule, theme.Rule)
        text_x += HexSpace * 2

        // text value
        o.print(text_x, line_y, line.str, theme.Output)

        // mark found positions
        // for c, f := range heap.Chain {
        //     o.hexMark(hex_x, line_y, c, line.str, f.Name)
        //     o.textMark(text_x, line_y, c, line.str, f.Name)
        // }
    }
}

func (o *Output) hexBuffer(heap *heap.Heap, w, h int) (hl []hexLine, my int) {
    c := int(float64((w - (8 + HexSpace)) + HexSpace) / 3.5)
    c -= c % 2

    hw := int(float64(c) * 2.5)

    my = len(heap.MMap) / c

    if len(heap.MMap) % c > 0 {
        my++
    }

    m := heap.MMap[o.delta_y * c:]

    for i := 0; i < len(m); i += c {
        if len(hl) >= h {
            return hl[:h], my
        }

        l := hexLine{
            off: fmt.Sprintf("%0*X ", 8, o.delta_y + i),
            hex: "",
            str: "",
        }

        for j := 0; j < c; j++ {
            if i + j >= len(m) {
                break
            }

            b := m[i + j]

            l.str = fmt.Sprintf("%s%c", l.str, b)
            l.hex = fmt.Sprintf("%s%02X", l.hex, b)

            // make hex gap
            if (j+1) % 2 == 0 {
                l.hex += " "
            }
        }

        l.hex = fmt.Sprintf("%-*s", hw, l.hex)

        hl = append(hl, l)
    }

    return
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

            o.print(x + dx, y, h, theme.Colors[c % len(theme.Colors)])
        }

        j = i+1
    }

    return
}

func (o *Output) hexGoto(s string) {
    // TODO
}
