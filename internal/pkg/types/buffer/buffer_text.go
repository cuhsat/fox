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
	Lines []TextLine
	Parts []TextPart

	SMap *smap.SMap
}

type TextLine struct {
	Line
}

type TextPart struct {
	Part
}

func (tb TextBuffer) String() string {
	var sb strings.Builder

	for i, l := range tb.Lines {
		sb.WriteString(l.String())

		if i < len(tb.Lines)-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

func (tl TextLine) String() string {
	return tl.Str
}

func Text(ctx *Context) TextBuffer {
	var tb TextBuffer

	d := text.Dec(ctx.Heap.Total())

	if ctx.Line {
		ctx.W -= (d + 1)
	}

	if ctx.Wrap && ctx.Heap.RMap() == nil {
		ctx.Heap.Wrap(ctx.W)
	}

	if ctx.Wrap {
		tb.SMap = ctx.Heap.RMap()
	} else {
		tb.SMap = ctx.Heap.SMap()
	}

	tb.W, tb.H = tb.SMap.Size()

	addLines(ctx, &tb, d)

	if len(tb.Lines) >= ctx.H {
		tb.Lines = tb.Lines[:ctx.H]
	}

	for _, f := range *types.Filters() {
		addParts(ctx, &tb, f)
	}

	return tb
}

func addLines(ctx *Context, tb *TextBuffer, d int) {
	var nr, str string

	for i, s := range (*tb.SMap)[ctx.Y:] {
		if i >= ctx.H {
			break
		}

		nr = fmt.Sprintf("%0*d", d, s.Nr)
		str = trim(format(ctx.Heap.Unmap(&s), &s), ctx.X, ctx.W)

		tb.Lines = append(tb.Lines, TextLine{
			Line: Line{Nr: nr, Str: str},
		})
	}
}

func addParts(ctx *Context, tb *TextBuffer, f string) {
	var str string
	var x int

	re, _ := regexp.Compile(f)

	for y, s := range tb.Lines {
		if y >= ctx.H {
			break
		}

		for _, i := range re.FindAllStringIndex(s.Str, -1) {
			x = text.Len(s.Str[:i[0]])
			str = s.Str[i[0]:i[1]]

			tb.Parts = append(tb.Parts, TextPart{
				Part: Part{X: x, Y: y, Str: str},
			})
		}
	}
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
