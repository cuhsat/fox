package widget

import (
    "fmt"
    "math"
    "strconv"
    "strings"

    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/ui/theme"
)

const (
    LineSpace = 1
)

type lineData struct {
    nr, str string
}

func (o *Output) lineRender(heap *heap.Heap, x, y, w, h int) {
    // convert logical to display lines
    lines := o.lineBuffer(heap, w, h)

    if len(lines) > 0 {
        w -= len(lines[0].nr) + LineSpace
    }

    // set buffer bounds
    o.last_x = max(heap.SMap.Width() - w, 0)
    o.last_y = max(len(heap.SMap) - h, 0)

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
        for c, f := range heap.Chain {
            o.lineMark(line_x, line_y, c, line.str, f.Name)
        }
    }
}

func (o *Output) lineBuffer(heap *heap.Heap, w, h int) (ld []lineData) {
    len_nr := int(math.Log10(float64(heap.Length()))) + 1

    if o.status.Line {
        w -= (len_nr + LineSpace)
    }

    for i, se := range heap.SMap[o.delta_y:] {
        if len(ld) >= h {
            return ld[:h]
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
        if o.status.Wrap {
            for {
                if length(str) < w+1 {
                    break
                }

                ld = append(ld, lineData{
                    nr: nr,
                    str: str[:w-1] + "\r",
                })

                str = str[w-1:]
            }
        }

        ld = append(ld, lineData{
            nr: nr,
            str: str,
        })
    }

    return
}

func (o *Output) lineMark(x, y, c int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    o.print(x + i, y, f, theme.Colors[c % len(theme.Colors)])
    
    o.lineMark(x + i+1, y, c, s[i+1:], f)
}

func (o *Output) lineGoto(s string) {
    switch s[0] {
    case '+':
        delta, _ := strconv.Atoi(s[1:])
        o.ScrollDown(delta)

    case '-':
        delta, _ := strconv.Atoi(s[1:])
        o.ScrollUp(delta)

    default:
        delta, _ := strconv.Atoi(s)
        o.delta_y = max(min(delta-1, o.last_y), 0)
    }
}
