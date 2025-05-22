package buffer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

const (
	BlankTab = "    "
	BlankOff = " "
)

type TextBuffer struct {
	Buffer
	Width int

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
	buf.Width = text.Dec(ctx.Heap.Total())

	if ctx.Line {
		ctx.W -= buf.Width + 1
	}

	if ctx.Wrap && ctx.Heap.RMap() == nil {
		ctx.Heap.Wrap(ctx.W) // TODO: ctx.W not wrapped correctly
	}

	if ctx.Wrap {
		buf.SMap = ctx.Heap.RMap()
	} else {
		buf.SMap = ctx.Heap.SMap()
	}

	buf.W, buf.H = buf.SMap.Size()

	buf.Lines = make(chan TextLine, Size)
	buf.Parts = make(chan TextPart, Size)

	var re *regexp.Regexp
	var fs = *types.GetFilters()

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

			s := trim(format(ctx.Heap.Unmap(&str), &str), ctx.X, ctx.W)

			buf.Lines <- TextLine{
				Line: Line{
					Nr:  fmt.Sprintf("%0*d", buf.Width, str.Nr),
					Str: s,
				},
			}

			if re == nil {
				continue
			}

			for _, i := range re.FindAllStringIndex(s, -1) {
				buf.Parts <- TextPart{
					Part: Part{
						X:   text.Len(s[:i[0]]),
						Y:   y,
						Str: s[i[0]:i[1]],
					},
				}
			}
		}
	}()

	return
}

func format(s string, str *smap.String) string {
	// replace tabulators
	s = strings.ReplaceAll(s, "\t", BlankTab)

	// prepend blank for offset
	if str.Off > 0 {
		s = strings.Repeat(BlankOff, int(str.Off)) + strings.TrimSpace(s)
	}

	return s
}

func trim(s string, x, w int) string {
	return text.Trim(s, min(x, text.Len(s)), w)
}
