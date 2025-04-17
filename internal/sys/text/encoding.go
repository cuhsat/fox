package text

import (
    "strings"
    "unicode"
)

const (
    minASCII = 0x20
    maxASCII = 0x7f
)

const (
    notASCII = '.'
    notUnicode = '·'
)

func ToASCII(s string) string {
    var sb strings.Builder

    for _, r := range s {
        sb.WriteRune(AsASCII(r))
    }

    return sb.String()
}

func AsASCII(r rune) rune {
    if r < minASCII || r > maxASCII {
        return notASCII
    } else {
        return r
    }
}

func AsUnicode(r rune) rune {
    if !unicode.IsPrint(r) {
        return notUnicode
    } else {
        return r
    }
}
