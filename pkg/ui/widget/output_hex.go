package widget

import (
    "fmt"
    "strings"

    "github.com/cuhsat/cu/pkg/ui/theme"
)

const (
    HexSpace = 1
)

const (
    Rule = "┃"
)

type hexData struct {
    off, hex, str string
}

func (o *Output) hexRender(x, y, w, h int) {
    // convert logical to display lines
    lines, max_y := o.hexBuffer(w, h)

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

        hex_x := line_x

        if o.status.Line {
            // offset number
            o.print(hex_x, line_y, line.off, theme.Hint)
            hex_x += len(line.off)

            // offset separator
            o.print(hex_x, line_y, Rule, theme.Rule)
            hex_x += HexSpace * 2            
        }

        // hex values
        o.print(hex_x, line_y, line.hex, theme.Output)
        text_x := hex_x + len(line.hex)

        // hex separator
        o.print(text_x, line_y, Rule, theme.Rule)
        text_x += HexSpace * 2

        // text value
        o.printAscii(text_x, line_y, line.str, theme.Output)

        // mark found positions
        // for c, f := range heap.Chain {
        //     o.hexMark(hex_x, line_y, c, line.str, f.Name)
        //     o.textMark(text_x, line_y, c, line.str, f.Name)
        // }
    }
}

func (o *Output) hexBuffer(w, h int) (hd []hexData, my int) {
    l := 0

    if o.status.Line {
        l = 8
    }

    c := int(float64((w - (l + HexSpace)) + HexSpace) / 3.5)
    c -= c % 2

    hw := int(float64(c) * 2.5)

    my = len(o.heap.MMap) / c

    if len(o.heap.MMap) % c > 0 {
        my++
    }

    m := o.heap.MMap[o.delta_y * c:]

    for i := 0; i < len(m); i += c {
        if len(hd) >= h {
            return hd[:h], my
        }

        d := hexData{
            off: fmt.Sprintf("%0*X ", 8, o.delta_y + i),
            hex: "",
            str: "",
        }

        for j := 0; j < c; j++ {
            if i + j >= len(m) {
                break
            }

            b := m[i + j]

            d.str = fmt.Sprintf("%s%c", d.str, b)
            d.hex = fmt.Sprintf("%s%02X", d.hex, b)

            // make hex gap
            if (j+1) % 2 == 0 {
                d.hex += " "
            }
        }

        d.hex = fmt.Sprintf("%-*s", hw, d.hex)

        hd = append(hd, d)
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
