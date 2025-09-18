package page

import (
	"fmt"

	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/text"
)

type HexPage struct {
	Page
	Lines chan HexLine
}

type HexLine struct {
	Line
	Hex string
}

func (hl HexLine) String() string {
	// canonical form
	return fmt.Sprintf("%s %s|%-16s|", hl.Nr, hl.Hex, hl.Str)
}

func Hex(ctx *Context) (page HexPage) {
	var tail int

	mmap := *ctx.Heap.MMap()

	limit := flags.Get().Limits

	if limit.IsTail && limit.Bytes > 0 {
		tail = max(int(ctx.Heap.Len())-limit.Bytes, 0)
	}

	page.W, page.H = ctx.W, len(mmap)/16

	if len(mmap)%16 > 0 {
		page.H++
	}

	page.Lines = make(chan HexLine, Size)

	go func() {
		defer close(page.Lines)

		mmap = mmap[ctx.Y*16:]

		for i := 0; i < len(mmap); i += 16 {
			if i/16 >= ctx.H {
				return
			}

			nr := fmt.Sprintf("%08x ", tail+i+(ctx.Y*16))

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

				// make a hex gap every 8 bytes
				if (j+1)%8 == 0 {
					l.Hex += "  "
				} else {
					l.Hex += " "
				}
			}

			l.Hex = fmt.Sprintf("%-*s", 50, l.Hex)
			l.Str = text.ToASCII(l.Str)

			page.Lines <- l
		}
	}()

	return
}
