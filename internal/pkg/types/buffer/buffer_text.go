package buffer

import (
	"fmt"

	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type TextBuffer struct {
	Buffer
	Y int
	N int

	Lines chan TextLine
	Parts chan TextPart

	SMap *smap.SMap
}

type TextLine struct {
	Line
}

type TextPart struct {
	Part
}

func (tl TextLine) String() string {
	return tl.Str
}

func Text(ctx *Context) (buf TextBuffer) {
	buf.N = text.Dec(ctx.Heap.Count())

	if ctx.Line {
		ctx.W -= buf.N + 1
	}

	cache, key := ctx.Heap.Cache(), ctx.Hash()

	buf.SMap = cache[key]

	if buf.SMap == nil {
		buf.SMap = ctx.Heap.SMap()

		if ctx.Wrap && buf.SMap.CanIndent() {
			buf.SMap = buf.SMap.Indent()
		} else if ctx.Wrap {
			buf.SMap = buf.SMap.Wrap(ctx.W)
		} else {
			buf.SMap = buf.SMap.Render()
		}

		cache[key] = buf.SMap
	}

	buf.Y = ctx.Y

	if ctx.Nr > 0 {
		lastY := len(*buf.SMap) - ctx.H

		// find requested line
		buf.Y, _ = buf.SMap.Find(ctx.Nr)
		buf.Y = min(buf.Y, lastY)
	}

	buf.W, buf.H = buf.SMap.Size()

	buf.Lines = make(chan TextLine, Size)
	buf.Parts = make(chan TextPart, Size)

	go func() {
		defer close(buf.Lines)
		defer close(buf.Parts)

		re := ctx.Heap.LastFilter().Regex

		for y, str := range (*buf.SMap)[buf.Y:] {
			if y >= ctx.H {
				return
			}

			n := fmt.Sprintf("%0*d", buf.N, str.Nr)
			s := text.Trim(str.Str, min(ctx.X, text.Len(str.Str)), ctx.W)

			buf.Lines <- TextLine{Line{n, s}}

			if re == nil {
				continue
			}

			for _, i := range re.FindAllStringIndex(s, -1) {
				buf.Parts <- TextPart{Part{
					text.Len(s[:i[0]]),
					y,
					s[i[0]:i[1]],
				}}
			}
		}
	}()

	return
}
