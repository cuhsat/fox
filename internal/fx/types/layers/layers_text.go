package layers

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

type TextLayer struct {
    Layer
    Lines []TextLine
    SMap smap.SMap
}

type TextLine struct {
    Line
}

func (tl TextLayer) String() string {
    var sb strings.Builder

    for i, l := range tl.Lines {
        sb.WriteString(l.String())
        
        if i < len(tl.Lines)-1 {
            sb.WriteRune('\n')            
        }
    }

    return sb.String()
}

func (tl TextLine) String() string {
    return tl.Str
}

func Text(ctx *Context) []TextLayer {
    var tl TextLayer

    d := text.Dec(ctx.Heap.Length())

    tl.SMap = ctx.Heap.SMap

    if ctx.Line {
        ctx.W -= (d + TextSpace)
    }

    if ctx.Wrap && ctx.Heap.Fmt != nil {
        textFormat(ctx, &tl, d)
    } else {
        textNormal(ctx, &tl, d)
    }

    if len(tl.Lines) >= ctx.H {
        tl.Lines = tl.Lines[:ctx.H]
    }

    return []TextLayer{tl}
}

func textFormat(ctx *Context, tl *TextLayer, d int) {
    for _, s := range tl.SMap {
        nr := fmt.Sprintf("%0*d", d, s.Nr)

        str := string(ctx.Heap.MMap[s.Start:s.End])

        for _, l := range ctx.Heap.Fmt(str) {
            tl.Lines = append(tl.Lines, TextLine{
                Line: Line{Nr: nr, Str: textFit(l, ctx.X, ctx.W)},
            })

            tl.W, tl.H  = max(tl.W, text.Len(l)), len(tl.Lines)
        }
    }

    tl.Lines = tl.Lines[ctx.Y:]
}

func textNormal(ctx *Context, tl *TextLayer, d int) {
    if ctx.Wrap {
        tl.SMap = tl.SMap.Wrap(ctx.W)
    }

    tl.W, tl.H = tl.SMap.Size()

    for i, s := range tl.SMap[ctx.Y:] {
        if i >= ctx.H {
            break
        }

        nr := fmt.Sprintf("%0*d", d, s.Nr)

        str := string(ctx.Heap.MMap[s.Start:s.End])
        
        // replace tabulators
        str = strings.ReplaceAll(str, "\t", TabSpace)

        tl.Lines = append(tl.Lines, TextLine{
            Line: Line{Nr: nr, Str: textFit(str, ctx.X, ctx.W)},
        })
    }
}

func textFit(s string, x, w int) string {
    return text.Trim(s, min(x, text.Len(s)), w)
}
