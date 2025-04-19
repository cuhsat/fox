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

func Abr(s string, x, w int) string {
    return runewidth.Truncate(s, w-x, "…")
}

func Len(s string) (l int) {
    return runewidth.StringWidth(s)
}

func Title(s string, w int) string {
    return Block([]string{s}, w)
}

func Block(s []string, w int) (r string) {
    if w < 0 {
        for _, ss := range s {
            w = max(w, len(ss))
        }
        w += 4
    }

    l := strings.Repeat("─", w-2)

    // header
    r += fmt.Sprintf("┌%s┐\n", l)

    // body
    for _, ss := range s {
        r += fmt.Sprintf("│ %-*s │\n", w-4, ss)
    }

    // footer
    r += fmt.Sprintf("└%s┘", l)

    return
}
