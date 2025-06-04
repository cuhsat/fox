package buffer

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/text"
)

type CsvBuffer struct {
	Buffer
	N int

	Lines chan CsvLine
}

type CsvLine struct {
	Line
	Cells []string
}

func Csv(ctx *Context) (buf CsvBuffer) {
	buf.N = text.Dec(ctx.Heap.Count())
	buf.W, buf.H = ctx.W, ctx.H

	buf.Lines = make(chan CsvLine, Size)

	smap := ctx.Heap.SMap()

	go func() {
		defer close(buf.Lines)

		r := csv.NewReader(strings.NewReader(smap.String()))

		cols, err := r.ReadAll()

		if err != nil {
			// TODO
		}

		for y, rows := range cols[ctx.Y:] {
			if y >= ctx.H {
				return
			}

			n := fmt.Sprintf("%0*d", buf.N, y)

			buf.Lines <- CsvLine{
				Line{n, ""},
				rows,
			}
		}
	}()

	return
}
