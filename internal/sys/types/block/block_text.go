package block

import (
    "fmt"
    "math"
    "strings"

    "github.com/cuhsat/fx/internal/sys"
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

    if ctx.Wrap && ctx.Heap.Fmt == nil {
        tb.SMap = tb.SMap.Wrap(ctx.W)
    }

    if ctx.Heap.Fmt != nil {
        // TODO: set size
        for _, s := range tb.SMap[ctx.Y:] {
            nr := fmt.Sprintf("%0*d", len_nr, s.Nr)

            str := string(ctx.Heap.MMap[s.Start:s.End])

            sys.Debug(">>>", str)

            if len(str) == 0 {
                tb.Lines = append(tb.Lines, TextLine{
                    Line: Line{Nr: nr, Str: ""},
                })

                continue
            }

            for _, l := range ctx.Heap.Fmt(str) {
                if len(tb.Lines) >= ctx.H {
                    break
                }

                // str = str[min(ctx.X, text.Length(str)):]

                if len(str) > ctx.W {
                    str = str[:ctx.W-1] + "\r"
                }

                tb.Lines = append(tb.Lines, TextLine{
                    Line: Line{Nr: nr, Str: l},
                })
            }
        }

        tb.W, tb.H = 80, len(tb.Lines)

    } else {
        tb.W, tb.H = tb.SMap.Size()

        for _, s := range tb.SMap[ctx.Y:] {
            if len(tb.Lines) >= ctx.H {
                break
            }

            nr := fmt.Sprintf("%0*d", len_nr, s.Nr)

            str := string(ctx.Heap.MMap[s.Start:s.End])
            str = str[min(ctx.X, text.Length(str)):]

            if len(str) > ctx.W {
                str = str[:ctx.W-1] + "\r"
            }

            tb.Lines = append(tb.Lines, TextLine{
                Line: Line{Nr: nr, Str: str},
            })
        }
    }

    if len(tb.Lines) >= ctx.H {
        tb.Lines = tb.Lines[:ctx.H]
    }

    return
}
