package utils

import (
    "strings"
    "unicode"
)

const (
    MinASCII = 0x20
    MaxASCII = 0x7f
)

const (
    NotASCII = '.'
    NotUnicode = '·'
)

func ToASCII(s string) string {
    var sb strings.Builder

    for _, r := range s {
        sb.WriteRune(AsASCII(r))
    }

    return sb.String()
}

func AsASCII(r rune) rune {
    if r < MinASCII || r > MaxASCII {
        return NotASCII
    } else {
        return r
    }
}

func AsUnicode(r rune) rune {
    if !unicode.IsPrint(r) {
        return NotUnicode
    } else {
        return r
    }
}
