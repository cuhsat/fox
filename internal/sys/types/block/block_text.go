package block

import (
    "fmt"
    "math"
    "strings"

    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/cuhsat/fx/internal/sys/types/smap"
)

const (
    SpaceText = 1
)

type TextBlock struct {
    Lines []TextLine

    SMap smap.SMap

    Block
}

type TextLine struct {
    Line
}

func (tb TextBlock) String() string {
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

func Text(ctx Context) (tb TextBlock) {
    len_nr := int(math.Log10(float64(ctx.Heap.Length()))) + 1

    tb.SMap = ctx.Heap.SMap

    if ctx.Line {
        ctx.W -= (len_nr + SpaceText)
    }

    if ctx.Wrap && ctx.Heap.Fmt != nil {
        textFormat(ctx, len_nr, &tb)
    } else {
        textNormal(ctx, len_nr, &tb)
    }

    if len(tb.Lines) >= ctx.H {
        tb.Lines = tb.Lines[:ctx.H]
    }

    return
}

func textFormat(ctx Context, nrw int, tb *TextBlock) {
    for _, s := range tb.SMap {
        nr := fmt.Sprintf("%0*d", nrw, s.Nr)

        str := string(ctx.Heap.MMap[s.Start:s.End])

        for _, l := range ctx.Heap.Fmt(str) {
            tb.Lines = append(tb.Lines, TextLine{
                Line: Line{Nr: nr, Str: textFit(ctx, l)},
            })

            tb.W, tb.H  = max(tb.W, text.Len(l)), len(tb.Lines)
        }
    }

    tb.Lines = tb.Lines[ctx.Y:]
}

func textNormal(ctx Context, nrw int, tb *TextBlock) {
    if ctx.Wrap {
        tb.SMap = tb.SMap.Wrap(ctx.W)
    }

    tb.W, tb.H = tb.SMap.Size()

    for i, s := range tb.SMap[ctx.Y:] {
        if i >= ctx.H {
            break
        }

        nr := fmt.Sprintf("%0*d", nrw, s.Nr)

        str := string(ctx.Heap.MMap[s.Start:s.End])

        tb.Lines = append(tb.Lines, TextLine{
            Line: Line{Nr: nr, Str: textFit(ctx, str)},
        })
    }
}

func textFit(ctx Context, s string) string {
    s = text.Pos(s, min(ctx.X, text.Len(s)))

    if text.Len(s) > ctx.W {
        s = text.Cut(s, ctx.W) + "\r"
    }

    return s
}
