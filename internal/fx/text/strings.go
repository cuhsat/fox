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

func Trim(s string, l, r int) string {
    s = runewidth.TruncateLeft(s, l, "")
    s = runewidth.Truncate(s, r, "→")

    return s
}

func Title(s string, w int) string {
    return Block([]string{s}, w, "")
}

func Block(s []string, w int, t string) (r string) {
    if w < 0 {
        for _, ss := range s {
            w = max(w, len(ss))
        }
        
        w += 4
    }

    a := strings.Repeat("─", w-2)
    b := a

    // title
    if len(t) > 0 {
        a = fmt.Sprintf("─ %s %s", t, strings.Repeat("─", len(t)-1))
    }

    // header
    r += fmt.Sprintf("┌%s┐\n", a)

    // body
    for _, ss := range s {
        r += fmt.Sprintf("│ %-*s │\n", w-4, ss)
    }

    // footer
    r += fmt.Sprintf("└%s┘", b)

    return
}
