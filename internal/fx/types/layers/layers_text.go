package layers

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/smap"
)

const (
    TextSpace = 1
)

const (
    TabBlank = "    "
    OffBlank = " "
)

type TextLayer struct {
    Layer
    Lines []TextLine
    Parts []TextPart

    SMap smap.SMap
}

type TextLine struct {
    Line
}

type TextPart struct {
    Part
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

func Text(ctx *Context) TextLayer {
    var tl TextLayer

    d := text.Dec(ctx.Heap.Length())

    if ctx.Line {
        ctx.W -= (d + TextSpace)
    }

    if ctx.Wrap {
        ctx.Heap.Wrap(ctx.W)
    } else {
        ctx.Heap.Reset()
    }

    // prioritize render map
    if ctx.Heap.RMap != nil {
        tl.SMap = ctx.Heap.RMap
    } else {
        tl.SMap = ctx.Heap.SMap
    }

    tl.W, tl.H = tl.SMap.Size()

    addLines(ctx, &tl, d)

    if len(tl.Lines) >= ctx.H {
        tl.Lines = tl.Lines[:ctx.H]
    }

    for _, f := range *types.GetFilters() {
        addParts(ctx, &tl, f)
    }

    return tl
}

func addLines(ctx *Context, tl *TextLayer, d int) {
    for i, s := range tl.SMap[ctx.Y:] {
        if i >= ctx.H {
            break
        }

        nr := fmt.Sprintf("%0*d", d, s.Nr)

        str := unmap(ctx, s)

        tl.Lines = append(tl.Lines, TextLine{
            Line: Line{Nr: nr, Str: trim(str, ctx.X, ctx.W)},
        })
    }
}

func addParts(ctx *Context, tl *TextLayer, f string) {
    re, _ := regexp.Compile(f)

    for y, s := range tl.Lines {
        if y >= ctx.H {
            break
        }

        for _, i := range re.FindAllStringIndex(s.Str, -1) {
            x := text.Len(s.Str[:i[0]])
            t := s.Str[i[0]:i[1]]

            tl.Parts = append(tl.Parts, TextPart{
                Part: Part{X: x, Y: y, Str: t},
            })
        }
    }
}

func unmap(ctx *Context, s *smap.String) string {
    str := string(ctx.Heap.MMap[s.Start:s.End])

    // replace tabulators
    str = strings.ReplaceAll(str, "\t", TabBlank)

    // prepend blank for offset
    if s.Off > 0 {
        str = strings.Repeat(OffBlank, s.Off) + strings.TrimSpace(str)
    }

    return str
}

func trim(s string, x, w int) string {
    return text.Trim(s, min(x, text.Len(s)), w)
}
