package text

import (
	"strings"
)

const (
	BEL = 7
	LF  = 10
	CR  = 13
	ESC = 27
)

func Unescape(s string) string {
	var sb strings.Builder

	ch := make(chan rune)

	go func(s string) {
		for _, r := range s {
			ch <- r
		}

		close(ch)
	}(s)

	for r := range ch {
		switch r {
		case ESC:
			switch r = <-ch; r {
			case '[':
				for r := range ch {
					if r != ';' && r != '?' && (r < '0' || r > '9') {
						break
					}
				}
			case ']':
				if r = <-ch; r >= 0 && r <= '9' {
					for r := range ch {
						switch r {
						case BEL:
							break
						case ESC:
							<-ch
							break
						}
					}
				}
			case '(', ')', '%':
				<-ch
			}
		case CR:
			if r := <-ch; r != LF {
				sb.WriteRune(r)
			}
		default:
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
