package buffer

import (
	"fmt"
	"strings"

	"github.com/cuhsat/fx/pkg/fx/text"
	"github.com/cuhsat/fx/pkg/fx/types"
)

const (
	SpaceHex = 1
)

type HexBuffer struct {
	Buffer
	Lines []HexLine
}

type HexLine struct {
	Line
	Hex string
}

func (hb HexBuffer) String() string {
	var sb strings.Builder

	for i, l := range hb.Lines {
		sb.WriteString(l.String())

		if i < len(hb.Lines)-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

func (hl HexLine) String() string {
	return fmt.Sprintf("%s %s %s", hl.Nr, hl.Hex, hl.Str)
}

func Hex(ctx *Context) HexBuffer {
	var hb HexBuffer

	off_w := 8 + SpaceHex

	cols := int(float64((ctx.W-off_w)+SpaceHex) / 3.5)
	cols -= cols % 2

	mmap := *ctx.Heap.MMap()

	tail := types.Limits().Tail.Bytes

	hex_w := int(float64(cols) * 2.5)

	hb.W, hb.H = ctx.W, len(mmap)/cols

	if len(mmap)%cols > 0 {
		hb.H++
	}

	mmap = mmap[ctx.Y*cols:]

	for i := 0; i < len(mmap); i += cols {
		if len(hb.Lines) >= ctx.H {
			hb.Lines = hb.Lines[:ctx.H]
			return hb
		}

		nr := fmt.Sprintf("%0*X ", 8, tail+ctx.Y+i)

		l := HexLine{
			Line: Line{Nr: nr, Str: ""},
			Hex:  "",
		}

		for j := range cols {
			if i+j >= len(mmap) {
				break
			}

			b := mmap[i+j]

			l.Hex = fmt.Sprintf("%s%02X", l.Hex, b)
			l.Str = fmt.Sprintf("%s%c", l.Str, b)

			// make hex gap
			if (j+1)%2 == 0 {
				l.Hex += " "
			}
		}

		l.Hex = fmt.Sprintf("%-*s", hex_w, l.Hex)
		l.Str = text.ToASCII(l.Str)

		hb.Lines = append(hb.Lines, l)
	}

	return hb
}
