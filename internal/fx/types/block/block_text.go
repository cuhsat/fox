package block

import (
    "fmt"
    "strings"

    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types/smap"
)

const (
    TextSpace = 1
)

const (
    TabSpace = "    "
)

type TextBlock struct {
    Block
    Lines []TextLine
    SMap smap.SMap
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

func Text(ctx *Context) (tb TextBlock) {
    d := text.Dec(ctx.Heap.Length())

    tb.SMap = ctx.Heap.SMap

    if ctx.Line {
        ctx.W -= (d + TextSpace)
    }

    if ctx.Wrap && ctx.Heap.Fmt != nil {
        textFormat(ctx, &tb, d)
    } else {
        textNormal(ctx, &tb, d)
    }

    if len(tb.Lines) >= ctx.H {
        tb.Lines = tb.Lines[:ctx.H]
    }

    return
}

func textFormat(ctx *Context, tb *TextBlock, d int) {
    for _, s := range tb.SMap {
        nr := fmt.Sprintf("%0*d", d, s.Nr)

        str := string(ctx.Heap.MMap[s.Start:s.End])

        for _, l := range ctx.Heap.Fmt(str) {
            tb.Lines = append(tb.Lines, TextLine{
                Line: Line{Nr: nr, Str: textFit(l, ctx.X, ctx.W)},
            })

            tb.W, tb.H  = max(tb.W, text.Len(l)), len(tb.Lines)
        }
    }

    tb.Lines = tb.Lines[ctx.Y:]
}

func textNormal(ctx *Context, tb *TextBlock, d int) {
    if ctx.Wrap {
        tb.SMap = tb.SMap.Wrap(ctx.W)
    }

    tb.W, tb.H = tb.SMap.Size()

    for i, s := range tb.SMap[ctx.Y:] {
        if i >= ctx.H {
            break
        }

        nr := fmt.Sprintf("%0*d", d, s.Nr)

        str := string(ctx.Heap.MMap[s.Start:s.End])
        
        // replace tabulators
        str = strings.ReplaceAll(str, "\t", TabSpace)

        tb.Lines = append(tb.Lines, TextLine{
            Line: Line{Nr: nr, Str: textFit(str, ctx.X, ctx.W)},
        })
    }
}

func textFit(s string, x, w int) string {
    return text.Trim(s, min(x, text.Len(s)), w)
}
