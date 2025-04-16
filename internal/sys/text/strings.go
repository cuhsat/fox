package text

import (
    "fmt"
    "math"
    "strings"

    "github.com/mattn/go-runewidth"
)

const (
    Dots = "…"
)

func Dec(n int) int {
    return int(math.Log10(float64(n))) + 1
}

func Pos(s string, x int) string {
    if x < Len(s) {
        return string([]rune(s)[x:])
    }

    return ""
}

func Cut(s string, w int) string {
    if w < Len(s) {
        return string([]rune(s)[:w-1])
    }

    return ""
}

func Abr(s string, x, w int) string {
    if x + Len(s) > w + 1 {
        s = Cut(s, (w-x)) + Dots
    }

    return s
}

func Len(s string) (l int) {
    for _, r := range s {
        l += runewidth.RuneWidth(r)
    }

    return
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
