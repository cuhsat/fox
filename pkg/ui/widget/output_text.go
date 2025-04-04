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
    TextSpace = 1
)

type textLine struct {
    nr, str string
}

func (o *Output) textRender(heap *heap.Heap, x, y, w, h int) {
    // convert logical to display lines
    lines := o.textBuffer(heap, w, h)

    if len(lines) > 0 {
        w -= len(lines[0].nr) + TextSpace
    }

    // set buffer bounds
    o.last_x = max(heap.SMap.Length() - w, 0)
    o.last_y = max(len(heap.SMap) - h, 0)

    // render buffer
    for i, line := range lines {
        line_x := x
        line_y := y + i

        // line number
        if o.line {
            o.print(line_x, line_y, line.nr, theme.Hint)
            line_x += len(line.nr) + TextSpace
        }

        // text value
        o.print(line_x, line_y, line.str, theme.Output)

        // mark found positions
        for c, f := range heap.Chain {
            o.textMark(line_x, line_y, c, line.str, f.Name)
        }
    }
}

func (o *Output) textBuffer(heap *heap.Heap, w, h int) (tl []textLine) {
    len_nr := int(math.Log10(float64(heap.Length()))) + 1

    if o.line {
        w -= (len_nr + TextSpace)
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

func (o *Output) textMark(x, y, c int, s, f string) {
    i := strings.Index(s, f)

    if i == -1 {
        return
    }

    o.print(x + i, y, f, theme.Colors[c % len(theme.Colors)])
    
    o.textMark(x + i+1, y, c, s[i+1:], f)
}

func (o *Output) textGoto(s string) {
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
