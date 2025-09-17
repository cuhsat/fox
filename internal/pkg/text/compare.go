package text

import (
	"fmt"
	"strings"

	"github.com/atombender/go-diff"
)

func Diff(a, b []string, l bool) string {
	var sb strings.Builder

	hunks := diff.Diff(a, b)

	n := Dec(max(len(a), len(b)))

	for _, h := range hunks {
		var r rune

		switch h.Operation {
		case diff.OpInsert:
			r = '+'
		case diff.OpDelete:
			r = '-'
		case diff.OpUnchanged:
			r = ' '
		}

		nr := fmt.Sprintf("%0*d", n, h.LineNum+1)

		if l {
			sb.WriteString(fmt.Sprintf("%c %s %s\n", r, nr, h.Line))
		} else {
			sb.WriteString(fmt.Sprintf("%c %s\n", r, h.Line))
		}
	}

	return sb.String()
}
