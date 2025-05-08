package text

import (
	"fmt"
	"math"
	"strings"

	"github.com/mattn/go-runewidth"
)

func Dec(n int) int {
	return int(math.Log10(float64(n)))+1
}

func Len(s string) (l int) {
	return runewidth.StringWidth(s)
}

func Abr(s string, w int) string {
	return runewidth.Truncate(s, w, "…")
}

func Abl(s string, w int) string {
	if Len(s) > w {
		s = "…" + runewidth.TruncateLeft(s, (Len(s)-w)+1, "")
	}

	return s
}

func Trim(s string, l, r int) string {
	s = runewidth.TruncateLeft(s, l, "")
	s = runewidth.Truncate(s, r, "→")

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
