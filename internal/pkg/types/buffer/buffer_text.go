package buffer

import (
	"fmt"
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type TextBuffer struct {
	Buffer
	Count int

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
	buf.Count = text.Dec(ctx.Heap.Count())

	if ctx.Line {
		ctx.W -= buf.Count + 1
	}

	fs := *types.GetFilters()
	id := hash(ctx, fs)
	ok := false

	// error!

	if buf.SMap, ok = Cache[id]; !ok {
		smap := ctx.Heap.SMap()

		if ctx.Wrap && smap.CanIndent() {
			buf.SMap = smap.Indent()
		} else if ctx.Wrap {
			buf.SMap = smap.Wrap(ctx.W)
		} else {
			buf.SMap = smap.Render()
		}

		Cache[id] = buf.SMap
	}

	buf.W, buf.H = buf.SMap.Size()

	buf.Lines = make(chan TextLine, Size)
	buf.Parts = make(chan TextPart, Size)

	var re *regexp.Regexp

	if len(fs) > 0 {
		re = regexp.MustCompile(fs[len(fs)-1])
	}

	go func() {
		defer close(buf.Lines)
		defer close(buf.Parts)

		for y, str := range (*buf.SMap)[ctx.Y:] {
			if y >= ctx.H {
				return
			}

			n := fmt.Sprintf("%0*d", buf.Count, str.Nr)
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

func hash(ctx *Context, fs types.Filters) string {
	var f string

	if len(fs) > 0 {
		f = fs[len(fs)-1]
	}

	return fmt.Sprintf("%s-%t-%t-%d-%d-%s",
		ctx.Heap.Path,
		ctx.Wrap,
		ctx.Line,
		ctx.W,
		ctx.H,
		f,
	)
}
