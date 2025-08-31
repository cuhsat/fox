package text

import (
	"fmt"
	"math"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/mattn/go-runewidth"
)

// Unicode icons
var UnicodeIcons = Icon{
	None: '·',
	HSep: '—',
	VSep: '∣',
	Size: '×',
	Grep: '❯',
	Ps1:  '❯',
}

// Default icons
var DefaultIcons = Icon{
	None: '·',
	HSep: '-',
	VSep: '|',
	Size: 'x',
	Grep: '>',
	Ps1:  '>',
}

type Icon struct {
	None, HSep, VSep, Size, Grep, Ps1 rune
}

func Dec(n int) int {
	return int(math.Log10(float64(n))) + 1
}

func Len(s string) (l int) {
	return runewidth.StringWidth(s)
}

func Abr(s string, w int) string {
	if Len(s) > w {
		s = runewidth.TruncateLeft(s, Len(s)-w, "…")
	}

	return s
}

func Pad(s string, w int) string {
	return runewidth.FillRight(s, w)
}

func Trim(s string, l, r int) string {
	s = runewidth.TruncateLeft(s, l, "")
	s = runewidth.Truncate(s, r, "")

	return s
}

func Title(s string, w int) (r string) {
	if w < 0 {
		w = 4 + len(s)
	}

	l := strings.Repeat("─", w-2)

	r += fmt.Sprintf("┌%s┐\n", l)
	r += fmt.Sprintf("│ %-*s │\n", w-4, s)
	r += fmt.Sprintf("└%s┘", l)

	return
}

func Icons() *Icon {
	if !flags.Get().UI.Legacy {
		return &UnicodeIcons
	} else {
		return &DefaultIcons
	}
}
