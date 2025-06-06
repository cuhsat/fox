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
	sr := strings.NewReader(ctx.Heap.SMap().String())
	cr := csv.NewReader(sr)

	cols, err := cr.ReadAll()

	if err != nil {
		sys.Error(err)
		return
	}

	buf.N = text.Dec(ctx.Heap.Count())
	buf.L = make([]int, len(cols[0]))

	// calculate row max length
	for _, rows := range cols {
		for i, row := range rows {
			buf.L[i] = max(text.Len(row), buf.L[i])
		}
	}

	// calculate buffer width
	for _, l := range buf.L {
		buf.W += l + 3
	}

	buf.W -= 2
	buf.H = len(cols)

	csvLine := func(nr int, ss []string) CsvLine {
		var sb strings.Builder

		for i, s := range ss {
			sb.WriteString(text.Padd(s, buf.L[i]))
			sb.WriteString("   ")
		}

		s := sb.String()

		return CsvLine{Line{
			fmt.Sprintf("%0*d", buf.N, nr),
			text.Trim(s, min(ctx.X, text.Len(s)), ctx.W),
		}}
	}

	buf.Lines = make(chan CsvLine, Size)

	go func() {
		defer close(buf.Lines)

		buf.Lines <- csvLine(0, cols[0])

		for y, rows := range cols[ctx.Y+1:] {
			if y >= ctx.H-1 {
				return
			}

			buf.Lines <- csvLine(ctx.Y+1+y, rows)
		}
	}()

	return
}
