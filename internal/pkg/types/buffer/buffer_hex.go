package buffer

import (
	"fmt"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
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
	// canonical form
	return fmt.Sprintf("%s %s|%-16s|", hl.Nr, hl.Hex, hl.Str)
}

func Hex(ctx *Context) HexBuffer {
	var hb HexBuffer
	var tail int

	mmap := *ctx.Heap.MMap()

	if types.GetLimits().Tail.Bytes > 0 {
		tail = max(int(ctx.Heap.Size())-types.GetLimits().Tail.Bytes, 0)
	}

	hb.W, hb.H = ctx.W, len(mmap)/16

	if len(mmap)%16 > 0 {
		hb.H++
	}

	mmap = mmap[ctx.Y*16:]

	for i := 0; i < len(mmap); i += 16 {
		if len(hb.Lines) >= ctx.H {
			hb.Lines = hb.Lines[:ctx.H]
			return hb
		}

		nr := fmt.Sprintf("%0*x ", 8, tail+i+(ctx.Y*16))

		l := HexLine{
			Line: Line{Nr: nr, Str: ""},
			Hex:  "",
		}

		for j := range 16 {
			if i+j >= len(mmap) {
				break
			}

			b := mmap[i+j]

			l.Hex = fmt.Sprintf("%s%02x", l.Hex, b)
			l.Str = fmt.Sprintf("%s%c", l.Str, b)

			// make hex gap every 8 bytes
			if (j+1)%8 == 0 {
				l.Hex += "  "
			} else {
				l.Hex += " "
			}
		}

		l.Hex = fmt.Sprintf("%-*s", 50, l.Hex)
		l.Str = text.ToASCII(l.Str)

		hb.Lines = append(hb.Lines, l)
	}

	return hb
}
