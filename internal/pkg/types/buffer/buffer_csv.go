package buffer

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
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
	SMap *smap.SMap

	Lens *[]int

	W int
	H int
}

func Csv(ctx *Context) (buf CsvBuffer) {
	buf.N = text.Dec(ctx.Heap.Count())

	cache, key := ctx.Heap.Cache(), ctx.Hash("csv")

	if val, ok := cache[key]; ok {
		buf.CMap = val.(*CMap)
	} else {
		var sb strings.Builder

		sr := strings.NewReader(ctx.Heap.SMap().String())
		cr := csv.NewReader(sr)

		cols, err := cr.ReadAll()

		if err != nil {
			panic(err) // TODO
		}

		smap := make(smap.SMap, len(cols))
		lens := make([]int, len(cols[0]))

		buf.CMap = &CMap{&smap, &lens, 0, 0}

		// calculate row max length
		for _, rows := range cols {
			for i, row := range rows {
				lens[i] = max(text.Len(row), lens[i])
			}
		}

		// calculate buffer width
		for _, l := range lens {
			buf.CMap.W += l + 3
		}

		// prepadd all rows
		for i, rows := range cols {
			for j, row := range rows {
				sb.WriteString(text.Padd(row, lens[j]))
				sb.WriteString("   ")
			}

			smap[i].Nr = i
			smap[i].Str = sb.String()

			sb.Reset()
		}

		buf.CMap.W -= 2
		buf.CMap.H = len(smap)

		cache[key] = buf.CMap
	}

	buf.W = buf.CMap.W
	buf.H = buf.CMap.H

	buf.Lines = make(chan CsvLine, Size)

	csvLine := func(str smap.String) CsvLine {
		return CsvLine{Line{
			fmt.Sprintf("%0*d", buf.N, str.Nr),
			text.Trim(str.Str, min(ctx.X, text.Len(str.Str)), ctx.W),
		}}
	}

	go func() {
		defer close(buf.Lines)

		buf.Lines <- csvLine((*buf.CMap.SMap)[0])

		for y, str := range (*buf.CMap.SMap)[ctx.Y+1:] {
			if y >= ctx.H-1 {
				return
			}

			buf.Lines <- csvLine(str)
		}
	}()

	return
}
