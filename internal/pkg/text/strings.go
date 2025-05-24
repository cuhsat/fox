package text

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/cuhsat/fox/internal/pkg/data"
)

const (
	Chevron = "❯"
)

func Dec(n int) int {
	return int(math.Log10(float64(n))) + 1
}

func Len(s string) (l int) {
	return runewidth.StringWidth(s)
}

func Abr(s string, w int) string {
	return runewidth.Truncate(s, w, "…")
}

func Abl(s string, w int) string {
	if Len(s) > w {
		s = runewidth.TruncateLeft(s, Len(s)-w, "…")
	}

	return s
}

func Trim(s string, l, r int) string {
	//s = runewidth.TruncateLeft(s, l, "")
	s = s[l:]

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

func Split(s string) (r []string) {
	n := 0

	for _, s := range strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case '"':
			n += 1
		case ' ':
			return n%2 == 0
		}

		return false
	}) {
		r = append(r, strings.ReplaceAll(s, "\"", ""))
	}

	return
}

func Reverse(s string) <-chan string {
	ch := make(chan string)

	go func() {
		for _, h := range data.Hashes {
			re := regexp.MustCompile(h.Regex)

			if re.MatchString(s) {
				for _, s := range h.Names {
					ch <- s
				}
			}
		}

		close(ch)
	}()

	return ch
}
