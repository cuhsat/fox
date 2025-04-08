package widget

import (
    "fmt"
    "math"
    "strconv"
    "strings"

    "github.com/cuhsat/cu/pkg/ui/theme"
)

const (
    LineSpace = 1
)

type textData struct {
    nr, str string
}

func (o *Output) textRender(x, y, w, h int) {
    lines, bw, bh := o.textBuffer(w, h)

    if len(lines) > 0 {
        w -= len(lines[0].nr) + LineSpace
    }

    // set buffer bounds
    o.last_x = max(bw - w, 0)
    o.last_y = max(bh - h, 0)

    // render buffer
    for i, line := range lines {
        line_x := x
        line_y := y + i

        // line number
        if o.status.Line {
            o.print(line_x, line_y, line.nr, theme.Hint)
            line_x += len(line.nr) + LineSpace
        }

        // text value
        o.print(line_x, line_y, line.str, theme.Output)

        // mark found positions
        for c, f := range o.heap.Chain {
            o.textMark(line_x, line_y, c, line.str, f.Name)
        }
    }
}

func (o *Output) textBuffer(w, h int) (td []textData, bw, bh int) {
    len_nr := int(math.Log10(float64(o.heap.Length()))) + 1

    o.smap = o.heap.SMap

    if o.status.Line {
        w -= (len_nr + LineSpace)
    }

    if o.status.Wrap {
        o.smap = o.smap.Wrap(w)
    }    

    bw, bh = o.smap.Size()

    for i, s := range o.smap[o.delta_y:] {
        if i >= h {
            break
        }

        nr := fmt.Sprintf("%0*d", len_nr, s.Nr)

        str := string(o.heap.MMap[s.Start:s.End])
        str = str[min(o.delta_x, length(str)):]

        if len(str) > w {
            str = str[:w-1] + "\r"
        }

        td = append(td, textData{
            nr: nr,
            str: str,
        })
    }

    if len(td) >= h {
        td = td[:h]
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

func (o *Output) textGoto(s string) {
    var nr int

    switch s[0] {
    case '+':
        delta, _ := strconv.Atoi(s[1:])
        nr = o.smap[o.delta_y].Nr + delta

    case '-':
        delta, _ := strconv.Atoi(s[1:])
        nr = o.smap[o.delta_y].Nr - delta

    default:
        nr, _ = strconv.Atoi(s)
    }

    y := o.smap.Find(nr)

    if y >= 0 {
        o.ScrollTo(o.delta_x, y)
    }
}
