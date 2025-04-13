package text

import (
    "fmt"
    "strings"

    "github.com/mattn/go-runewidth"
)

const (
    Cut = "…"
)

func Abbrev(s string, x, w int) string {
    if x + Length(s) > w + 1 {
        s = string([]rune(s)[:(w-x)-1]) + Cut
    }

    return s
}

func Length(s string) (l int) {
    for _, r := range s {
        l += runewidth.RuneWidth(r)
    }

    return
}

func Header(s string, w int) (h string) {
    b := [...]string{"┌","─","┐","│","└","┘"}
    
    l := strings.Repeat(b[1], w-2)

    h += fmt.Sprintf("%s%s%s\n", b[0], l, b[2])
    h += fmt.Sprintf("%s %-*s %s\n", b[3], w-4, s, b[3])
    h += fmt.Sprintf("%s%s%s", b[4], l, b[5])

    return
}
