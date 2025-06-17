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

	if ctx.Numbers {
		ctx.W -= buf.N + 1
	}

	cache, key := ctx.Heap.Cache(), ctx.Hash()

	if val, ok := cache[key]; ok {
		buf.SMap = val.(*smap.SMap)
	} else {
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
		lastY := max(len(*buf.SMap)-ctx.H, 0)

		// find the requested line
		buf.Y, _ = buf.SMap.Find(ctx.Nr)
		buf.Y = min(buf.Y, lastY)
	}

	buf.W, buf.H = buf.SMap.Size()

	buf.Lines = make(chan TextLine, Size)
	buf.Parts = make(chan TextPart, Size)

	go func() {
		defer close(buf.Lines)
		defer close(buf.Parts)

		fs := ctx.Heap.Filters()

		for y, str := range (*buf.SMap)[buf.Y:] {
			if y >= ctx.H {
				return
			}

			n := fmt.Sprintf("%0*d", buf.N, str.Nr)
			s := text.Trim(str.Str, min(ctx.X, text.Len(str.Str)), ctx.W)

			buf.Lines <- TextLine{Line{n, s}}

			for _, f := range fs {
				for _, i := range f.Regex.FindAllStringIndex(s, -1) {
					buf.Parts <- TextPart{Part{
						text.Len(s[:i[0]]),
						y,
						s[i[0]:i[1]],
					}}
				}
			}
		}
	}()

	return
}
