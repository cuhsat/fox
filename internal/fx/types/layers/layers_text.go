package layers

import (
    "fmt"
    "strings"

    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
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
    var tls []TextLayer
    var base TextLayer

    d := text.Dec(ctx.Heap.Length())

    base.SMap = ctx.Heap.SMap

    if ctx.Line {
        ctx.W -= (d + TextSpace)
    }

    if ctx.Wrap && ctx.Heap.Fmt != nil {
        textFormat(ctx, &base, d)
    } else {
        textNormal(ctx, &base, d)
    }

    if len(base.Lines) >= ctx.H {
        base.Lines = base.Lines[:ctx.H]
    }

    tls = append(tls, base)

    if !ctx.Wrap && ctx.Heap.Fmt == nil {
        for _, f := range *types.GetFilters() {
            tls = append(tls, textFilter(ctx, f))
        }
    }

    return tls
}

func textFormat(ctx *Context, tl *TextLayer, d int) {
    for _, s := range tl.SMap {
        nr := fmt.Sprintf("%0*d", d, s.Nr)

        str := unmap(ctx, s)

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

        str := unmap(ctx, s)

        tl.Lines = append(tl.Lines, TextLine{
            Line: Line{Nr: nr, Str: textFit(str, ctx.X, ctx.W)},
        })
    }
}

func textFilter(ctx *Context, f string) TextLayer {
    var tl TextLayer

    for i, s := range ctx.Heap.SMap[ctx.Y:] {
        if i >= ctx.H {
            break
        }

        str := unmap(ctx, s)

        // i, m := -1, ""

        // if ok, re := types.Regex(f); ok {
        //     l := re.FindIndex([]byte(s))

        //     if l != nil {
        //         i, m = l[0], s[l[0]:l[1]]
        //     }
        // } else {
        //     i, m = strings.Index(s, f), f
        // }

        // if i == -1 {
        //     return
        // }

        tl.Lines = append(tl.Lines, TextLine{
            Line: Line{Str: textFit(str, ctx.X, ctx.W)},
        })
    }

    return tl
}

func textFit(s string, x, w int) string {
    return text.Trim(s, min(x, text.Len(s)), w)
}

func unmap(ctx *Context, s *smap.String) string {
    str := string(ctx.Heap.MMap[s.Start:s.End])

    // replace tabulators
    str = strings.ReplaceAll(str, "\t", TabSpace)

    return str
}
