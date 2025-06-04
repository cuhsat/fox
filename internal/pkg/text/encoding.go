package text

import (
	"strings"
	"unicode"
)

const (
	LRE = '\u202a'
	RLE = '\u202b'
	LRO = '\u202d'
	RLO = '\u202e'
	LRI = '\u2066'
	RLI = '\u2067'
	FSI = '\u2068'
	PDF = '\u202c'
	PDI = '\u2069'
)

const (
	MinASCII = 0x20
	MaxASCII = 0x7f
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
		return '.'
	}

	return r
}

func AsUnicode(r rune) rune {
	// mitigate CVE-2021-42574
	switch r {
	case LRE, RLE, LRO, RLO, LRI, RLI, FSI, PDF, PDI:
		return '×'
	default:
		if !unicode.IsPrint(r) {
			return '·'
		}
	}

	return r
}
