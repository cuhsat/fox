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

	CMap *CMap
}

type CsvLine struct {
	Line
}

type CMap struct {
	Strs []string
	Rows []int

	w int
	h int
}

func Csv(ctx *Context) (buf CsvBuffer) {
	cache, key := ctx.Heap.Cache(), ctx.Hash("csv")

	if val, ok := cache[key]; ok {
		buf.CMap = val.(*CMap)
	} else {
		sr := strings.NewReader(ctx.Heap.SMap().String())
		cr := csv.NewReader(sr)

		cols, err := cr.ReadAll()

		if err != nil {
			panic(err) // TODO
		}

		buf.CMap = &CMap{
			make([]string, len(cols)),
			make([]int, len(cols[0])),
			0,
			0,
		}

		// calculate row max length
		for _, rows := range cols {
			for r, row := range rows {
				buf.CMap.Rows[r] = max(text.Len(row), buf.CMap.Rows[r])
			}
		}

		// calculate buffer width
		for _, l := range buf.CMap.Rows {
			buf.CMap.w += l + 3
		}

		// prepadd all rows
		for c, rows := range cols {
			var sb strings.Builder

			for r, row := range rows {
				sb.WriteString(text.Padd(row, buf.CMap.Rows[r]))
				sb.WriteString("   ")
			}

			buf.CMap.Strs[c] = sb.String()
		}

		buf.CMap.w -= 2
		buf.CMap.h = len(buf.CMap.Strs)

		cache[key] = buf.CMap
	}

	buf.W = buf.CMap.w
	buf.H = buf.CMap.h

	buf.N = text.Dec(ctx.Heap.Count())

	buf.Lines = make(chan CsvLine, Size)

	csvLine := func(nr int, s string) CsvLine {
		return CsvLine{Line{
			fmt.Sprintf("%0*d", buf.N, nr),
			text.Trim(s, min(ctx.X, text.Len(s)), ctx.W),
		}}
	}

	go func() {
		defer close(buf.Lines)

		buf.Lines <- csvLine(0, buf.CMap.Strs[0])

		for y, s := range buf.CMap.Strs[ctx.Y+1:] {
			if y >= ctx.H-1 {
				return
			}

			buf.Lines <- csvLine(ctx.Y+1+y, s)
		}
	}()

	return
}
