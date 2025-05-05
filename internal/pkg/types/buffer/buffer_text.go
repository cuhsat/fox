package buffer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/edsrzf/mmap-go"

	"github.com/cuhsat/fx/internal/pkg/text"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/smap"
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
	m := ctx.Heap.MMap()

	for i, s := range (*tb.SMap)[ctx.Y:] {
		if i >= ctx.H {
			break
		}

		nr := fmt.Sprintf("%0*d", d, s.Nr)

		str := unmap(m, &s)

		tb.Lines = append(tb.Lines, TextLine{
			Line: Line{Nr: nr, Str: trim(str, ctx.X, ctx.W)},
		})
	}
}

func addParts(ctx *Context, tb *TextBuffer, f string) {
	re, _ := regexp.Compile(f)

	for y, s := range tb.Lines {
		if y >= ctx.H {
			break
		}

		for _, i := range re.FindAllStringIndex(s.Str, -1) {
			x := text.Len(s.Str[:i[0]])
			t := s.Str[i[0]:i[1]]

			tb.Parts = append(tb.Parts, TextPart{
				Part: Part{X: x, Y: y, Str: t},
			})
		}
	}
}

func unmap(m *mmap.MMap, s *smap.String) string {
	str := string((*m)[s.Start:s.End])

	// replace tabulators
	str = strings.ReplaceAll(str, "\t", BlankTab)

	// prepend blank for offset
	if s.Off > 0 {
		str = strings.Repeat(BlankOff, int(s.Off)) + strings.TrimSpace(str)
	}

	return str
}

func trim(s string, x, w int) string {
	return text.Trim(s, min(x, text.Len(s)), w)
}
