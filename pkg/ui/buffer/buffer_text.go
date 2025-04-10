package buffer

import (
    "fmt"
    "math"
    "strings"

    "github.com/cuhsat/cu/pkg/fs/smap"
    "github.com/cuhsat/cu/pkg/fs/utils"
)

const (
    TextSpace = 1
)

type TextBuffer struct {
    Lines []TextLine

    SMap smap.SMap

    Buffer
}

type TextLine struct {
    Line
}

func Text(ctx Context) (tb TextBuffer) {
    len_nr := int(math.Log10(float64(ctx.Heap.Length()))) + 1

    tb.SMap = ctx.Heap.SMap

    if ctx.Line {
        ctx.W -= (len_nr + TextSpace)
    }

    if ctx.Wrap {
        tb.SMap = tb.SMap.Wrap(ctx.W)
    }    

    tb.W, tb.H = tb.SMap.Size()

    for i, s := range tb.SMap[ctx.Y:] {
        if i >= ctx.H {
            break
        }

        nr := fmt.Sprintf("%0*d", len_nr, s.Nr)

        str := string(ctx.Heap.MMap[s.Start:s.End])
        str = str[min(ctx.X, utils.Length(str)):]

        if len(str) > ctx.W {
            str = str[:ctx.W-1] + "\r"
        }

        tb.Lines = append(tb.Lines, TextLine{
            Line: Line{Nr: nr, Str: str},
        })
    }

    if len(tb.Lines) >= ctx.H {
        tb.Lines = tb.Lines[:ctx.H]
    }

    return
}

func (tb TextBuffer) String() string {
    var sb strings.Builder

    for i, l := range tb.Lines {
        sb.WriteString(l.String())
        
        if i < len(tb.Lines)-1 {
            sb.WriteRune('\n')            
        }
    }

    return sb.String()
}

func (tl TextLine) String() string {
    return tl.Str
}
