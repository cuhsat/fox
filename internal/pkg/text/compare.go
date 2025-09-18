package text

import (
	"fmt"
	"strings"

	"github.com/atombender/go-diff"

	"github.com/cuhsat/fox/internal/pkg/flags"
)

func Diff(f1, f2 string, l1, l2 []string, git bool) string {
	var sb strings.Builder

	flg := flags.Get()

	if git {
		sb.WriteString(fmt.Sprintf("--- %s\n", f1))
		sb.WriteString(fmt.Sprintf("+++ %s\n", f2))
	} else if !flg.NoFile {
		sb.WriteString(fmt.Sprintln(f1))
		sb.WriteString(fmt.Sprintln(f2))
	}

	n := Dec(max(len(l1), len(l2)))

	hunks := diff.Diff(l1, l2)

	for i, h := range hunks {
		var r rune

		switch h.Operation {
		case diff.OpInsert:
			r = '+'
		case diff.OpDelete:
			r = '-'
		case diff.OpUnchanged:
			r = ' '
			if git {
				continue
			}
		}

		nr := fmt.Sprintf("%0*d", n, h.LineNum+1)

		if git {
			if i < len(hunks)-1 && hunks[i].LineNum == hunks[i+1].LineNum {
				sb.WriteString(fmt.Sprintf("@@ -%d +%d @@\n", h.LineNum+1, h.LineNum+1))
				sb.WriteString(fmt.Sprintf("-%s\n", h.Line))
				sb.WriteString(fmt.Sprintf("+%s\n", hunks[i+1].Line))
			} else if i > 0 && hunks[i].LineNum == hunks[i-1].LineNum {
				continue // skip double changes
			} else {
				sb.WriteString(fmt.Sprintf("@@ %c%d @@\n", r, h.LineNum+1))
				sb.WriteString(fmt.Sprintf("%c%s\n", r, h.Line))
			}
		} else if !flg.NoLine {
			sb.WriteString(fmt.Sprintf("%c %s %s\n", r, nr, h.Line))
		} else {
			sb.WriteString(fmt.Sprintf("%c %s\n", r, h.Line))
		}
	}

	return sb.String()
}
