package layers

import (
    "fmt"
    "strings"

    "github.com/cuhsat/fx/internal/fx/text"
)

const (
    HexSpace = 1
)

type HexLayer struct {
    Layer
    Lines []HexLine
}

type HexLine struct {
    Line
    Hex string
}

func (hl HexLayer) String() string {
    var sb strings.Builder

    for i, l := range hl.Lines {
        sb.WriteString(l.String())
        
        if i < len(hl.Lines)-1 {
            sb.WriteRune('\n')            
        }
    }

    return sb.String()
}

func (hl HexLine) String() string {
    return fmt.Sprintf("%s %s %s", hl.Nr, hl.Hex, hl.Str)
}

func Hex(ctx *Context) HexLayer {
    var hl HexLayer

    off_w := 8 + HexSpace

    cols := int(float64((ctx.W - off_w) + HexSpace) / 3.5)
    cols -= cols % 2

    mmap := ctx.Heap.MMap
    tail := ctx.Heap.Tail

    hex_w := int(float64(cols) * 2.5)

    hl.W, hl.H = ctx.W, len(mmap) / cols

    if len(mmap) % cols > 0 {
        hl.H++
    }

    mmap = mmap[ctx.Y * cols:]

    for i := 0; i < len(mmap); i += cols {
        if len(hl.Lines) >= ctx.H {
            hl.Lines = hl.Lines[:ctx.H]
            return hl
        }

        nr := fmt.Sprintf("%0*X ", 8, tail + ctx.Y + i)

        l := HexLine{
            Line: Line{Nr: nr, Str: ""},
            Hex: "",
        }

        for j := 0; j < cols; j++ {
            if i + j >= len(mmap) {
                break
            }

            b := mmap[i + j]

            l.Hex = fmt.Sprintf("%s%02X", l.Hex, b)
            l.Str = fmt.Sprintf("%s%c", l.Str, b)

            // make hex gap
            if (j+1) % 2 == 0 {
                l.Hex += " "
            }
        }

        l.Hex = fmt.Sprintf("%-*s", hex_w, l.Hex)
        l.Str = text.ToASCII(l.Str)

        hl.Lines = append(hl.Lines, l)
    }

    return hl
}
