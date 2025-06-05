package buffer

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
)

type CsvBuffer struct {
	Buffer
	N int
	L []int

	Lines chan CsvLine
}

type CsvLine struct {
	Line
}

func Csv(ctx *Context) (buf CsvBuffer) {
	smap := ctx.Heap.SMap()

	r := csv.NewReader(strings.NewReader(smap.String()))

	cols, err := r.ReadAll()

	if err != nil {
		sys.Error(err)
		return
	}

	buf.N = text.Dec(ctx.Heap.Count())
	buf.L = make([]int, len(cols))
	buf.Lines = make(chan CsvLine, Size)

	// calculate cell max length
	for _, rows := range cols {
		for i, row := range rows {
			buf.L[i] = max(text.Len(row), buf.L[i])
		}
	}

	// calculate buffer width
	for _, w := range buf.L {
		buf.W += (w + 3)
	}

	buf.H = len(cols)

	go func() {
		defer close(buf.Lines)

		for y, rows := range cols[ctx.Y:] {
			var sb strings.Builder

			if y >= ctx.H {
				return
			}

			for i, row := range rows {
				sb.WriteString(text.PadR(row, buf.L[i]))
				sb.WriteString("   ")
			}

			n := fmt.Sprintf("%0*d", buf.N, (*smap)[ctx.Y+y].Nr)
			s := sb.String()

			s = text.Trim(s, min(ctx.X, text.Len(s)), ctx.W)

			buf.Lines <- CsvLine{
				Line{n, s},
			}
		}
	}()

	return
}
