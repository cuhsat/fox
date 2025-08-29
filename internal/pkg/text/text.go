package text

import (
	"fmt"
	"math"
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	PS1 = "❯"
)

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
