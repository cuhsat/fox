package page

import (
	"fmt"

	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

type TextPage struct {
	Page
	Y int
	N int

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

func Text(ctx *Context) (page TextPage) {
	page.N = text.Dec(ctx.Heap.Count())

	if ctx.Navi {
		ctx.W -= 2 + page.N
		ctx.H -= 1
	}

	key := ctx.Hash()

	if val, ok := ctx.Heap.Cache.Load(key); ok {
		page.FMap = val.(*smap.SMap)
	} else {
		page.FMap = ctx.Heap.FMap()

		if ctx.Wrap && page.FMap.CanFormat() {
			page.FMap = page.FMap.Format(ctx.Space)
		} else if ctx.Wrap {
			page.FMap = page.FMap.Wrap(ctx.Space, ctx.W)
		} else {
			page.FMap = page.FMap.Render(ctx.Space)
		}

		ctx.Heap.Cache.Store(key, page.FMap)
	}

	page.Y = ctx.Y

	if ctx.Nr > 0 {
		lastY := max(len(*page.FMap)-1, 0)

		// find the requested line
		page.Y, _ = page.FMap.Find(ctx.Nr)
		page.Y = min(page.Y, lastY)
	}

	page.W, page.H = page.FMap.Size()

	page.Lines = make(chan TextLine, Size)
	page.Parts = make(chan TextPart, Size)

	go func() {
		defer close(page.Lines)
		defer close(page.Parts)

		fs := ctx.Heap.Filters()

		sep, grp, num := 0, 0, 1

		// pinned head
		if ctx.Pinned {
			n := fmt.Sprintf("%0*d", page.N, 1)
			s := (*ctx.Heap.SMap())[0].Str
			s = text.Trim(s, min(ctx.X, text.Len(s)), ctx.W)

			page.Lines <- TextLine{Line{n, 0, s}}
		}

		for y, str := range (*page.FMap)[page.Y:] {
			if y >= ctx.H {
				return
			}

			// insert context separator
			if ctx.Context && grp != str.Grp && num > 1 {
				page.Lines <- TextLine{Line{"--", str.Grp, ""}}
				num = 1
				sep++
			}

			n := fmt.Sprintf("%0*d", page.N, str.Nr)
			s := text.Trim(str.Str, min(ctx.X, text.Len(str.Str)), ctx.W)

			page.Lines <- TextLine{Line{n, str.Grp, s}}

			if ctx.Pinned {
				y++
			}

			for _, f := range fs {
				for _, i := range f.Regex.FindAllStringIndex(s, -1) {
					page.Parts <- TextPart{Part{
						text.Len(s[:i[0]]),
						y + sep,
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
