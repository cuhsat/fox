package widget

import (
    "fmt"
    "math"
    "strings"

    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    Space = 1
)

type Output struct {
    widget

    hex, nrs, wrap bool
    last_x int
    last_y int
    delta_x int
    delta_y int
}

type textLine struct {
    nr, str string
}

type hexLine struct {
    off, hex, str string
} 

func NewOutput(screen tcell.Screen, hex bool) *Output {
    return &Output{
        widget: widget{
            screen: screen,
        },
        hex: hex,
        nrs: true,
        wrap: false,
        last_x: 0,
        last_y: 0,
        delta_x: 0,
        delta_y: 0,
    }
}

func (o *Output) Render(heap *data.Heap, x, y, w, h int) {
    if !o.hex {
        o.textRender(heap, x, y, w, h)
    } else {
        o.hexRender(heap, x, y, w, h)
    }
}

func (o *Output) Reset() {
    o.delta_x = 0
    o.delta_y = 0
}

func (o *Output) ScrollBegin() {
    o.delta_y = 0
}

func (o *Output) ScrollEnd() {
    o.delta_y = o.last_y
}

func (o *Output) ScrollUp(delta int) {
    o.delta_y = max(o.delta_y - delta, 0)
}

func (o *Output) ScrollDown(delta int) {
    o.delta_y = min(o.delta_y + delta, o.last_y)
}

func (o *Output) ScrollLeft(delta int) {
    o.delta_x = max(o.delta_x - delta, 0)
}

func (o *Output) ScrollRight(delta int) {
    o.delta_x = min(o.delta_x + delta, o.last_x)
}

func (o *Output) ScrollPageUp(delta int) {
    o.delta_y = max(o.delta_y - delta, 0)
}

func (o *Output) ScrollPageDown(delta int) {
    o.delta_y = min(o.delta_y + delta, o.last_y)
}

func (o *Output) ToggleNumbers() {
    o.nrs = !o.nrs
}

func (o *Output) ToggleWrap() {
    o.wrap = !o.wrap
}

func (o *Output) ToggleHex() {
    o.hex = !o.hex
    o.Reset()
}

func (o *Output) textRender(heap *data.Heap, x, y, w, h int) {
    // convert logical to display lines
    lines := o.textBuffer(heap, w, h)

    if len(lines) > 0 {
        w -= len(lines[0].nr) + Space
    }

    // set buffer bounds
    o.last_x = max(wid(heap.SMap) - w, 0)
    o.last_y = max(len(heap.SMap) - h, 0)

    // render buffer
    for i, line := range lines {
        line_x := x
        line_y := y + i

        // line number
        if o.nrs {
            o.print(line_x, line_y, line.nr, theme.Hint)
            line_x += len(line.nr) + Space
        }

        // text value
        o.print(line_x, line_y, line.str, theme.Output)

        // mark found positions
        for c, f := range heap.Chain {
            o.textMark(line_x, line_y, c, line.str, f.Name)
        }
    }
}

func (o *Output) textBuffer(heap *data.Heap, w, h int) (tl []textLine) {
    len_nr := int(math.Log10(float64(heap.Lines()))) + 1

    if o.nrs {
        w -= (len_nr + Space)
    }

    for i, se := range heap.SMap[o.delta_y:] {
        if len(tl) >= h {
            return tl[:h]
        }

        if i >= h {
            return
        }

        // line number
        nr := fmt.Sprintf("%0*d", len_nr, se.Nr)

        // logical line
        str := string(heap.MMap[se.Start:se.End])
        str = str[min(o.delta_x, length(str)):]

        // display lines
        if o.wrap {
            for {
                if length(str) < w+1 {
                    break
                }

                tl = append(tl, textLine{
                    nr: nr,
                    str: str[:w-1] + "\r",
                })

                str = str[w-1:]
            }
        }

        tl = append(tl, textLine{
            nr: nr,
            str: str,
        })
    }

    return
}

func (o *Output) hexRender(heap *data.Heap, x, y, w, h int) {
    // convert logical to display lines
    lines, max_y := o.hexBuffer(heap, w, h)

    if len(lines) > 0 {
        w -= len(lines[0].off) + Space
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
        hex_x := line_x + len(line.off) + Space

        // hex values
        o.print(hex_x, line_y, line.hex, theme.Output)
        text_x := hex_x + len(line.hex) + Space

        // text value
        o.print(text_x, line_y, line.str, theme.Output)

        // mark found positions
        for c, f := range heap.Chain {
            o.hexMark(hex_x, line_y, c, line.str, f.Name)
            o.textMark(text_x, line_y, c, line.str, f.Name)
        }
    }
}

func (o *Output) hexBuffer(heap *data.Heap, w, h int) (hl []hexLine, my int) {
    c := int(float64((w - (8 + Space)) + Space) / 3.5)
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
            l.hex = fmt.Sprintf("%s%02x", l.hex, b)

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

func (o *Output) textMark(x, y, c int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    o.print(x + i, y, f, theme.Colors[c % len(theme.Colors)])
    
    o.textMark(x + i+1, y, c, s[i+1:], f)
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

func wid(s data.SMap) (w int) {
    for _, se := range s {
        w = max(w, se.Len)
    }

    return
}
