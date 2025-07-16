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
	S int

	Lines chan TextLine
	Parts chan TextPart

	FMap *smap.SMap
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
		buf.FMap = val.(*smap.SMap)
	} else {
		buf.FMap = ctx.Heap.FMap()

		if ctx.Wrap && buf.FMap.CanIndent() {
			buf.FMap = buf.FMap.Indent()
		} else if ctx.Wrap {
			buf.FMap = buf.FMap.Wrap(ctx.W)
		} else {
			buf.FMap = buf.FMap.Render()
		}

		cache[key] = buf.FMap
	}

	buf.Y = ctx.Y

	if ctx.Nr > 0 {
		lastY := max(len(*buf.FMap)-1, 0)

		// find the requested line
		buf.Y, _ = buf.FMap.Find(ctx.Nr)
		buf.Y = min(buf.Y, lastY)
	}

	buf.W, buf.H = buf.FMap.Size()

	buf.Lines = make(chan TextLine, Size)
	buf.Parts = make(chan TextPart, Size)

	go func() {
		defer close(buf.Lines)
		defer close(buf.Parts)

		fs := ctx.Heap.Filters()

		grp, num := 0, 1

		for y, str := range (*buf.FMap)[buf.Y:] {
			if y >= ctx.H {
				return
			}

			// insert context separator
			if ctx.Context && grp != str.Grp && num > 1 {
				buf.Lines <- TextLine{Line{"--", str.Grp, ""}}
				buf.S++
				num = 1
			}

			n := fmt.Sprintf("%0*d", buf.N, str.Nr)
			s := text.Trim(str.Str, min(ctx.X, text.Len(str.Str)), ctx.W)

			buf.Lines <- TextLine{Line{n, str.Grp, s}}

			for _, f := range fs {
				for _, i := range f.Regex.FindAllStringIndex(s, -1) {
					buf.Parts <- TextPart{Part{
						text.Len(s[:i[0]]),
						y + buf.S,
						str.Grp,
						s[i[0]:i[1]],
					}}
				}
			}

			grp = str.Grp
			num++
		}
	}()

	return
}
