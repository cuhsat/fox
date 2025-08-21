package text

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cuhsat/fox/internal/pkg/flags"
)

type String struct {
	Off int
	Str string
}

func Carve(ch <-chan byte, str chan<- String, n, m int) {
	var rs []rune
	var off int

	flush := func() {
		if len(rs) >= n && len(rs) <= m {
			o := max(off-(len(rs)+1), 0)
			s := string(rs)

			if len(strings.TrimSpace(s)) > 0 {
				str <- String{o, s}
			}
		}

		rs = rs[:0]
	}

	defer close(str)
	defer flush()

	flg := flags.Get().Strings
	buf := make([]byte, 4)

	for b := range ch {
		buf[0] = b
		off++

		if flg.Ascii {
			if b >= MinASCII && b <= MaxASCII {
				rs = append(rs, rune(b))
			} else {
				flush()
			}
		} else {
			l := 1
			k := 1

			if b&0x80 == 0 {
				k = 1
			} else if b&0xE0 == 0xC0 {
				k = 2
			} else if b&0xF0 == 0xE0 {
				k = 3
			} else if b&0xF8 == 0xF0 {
				k = 4
			}

			if k > 1 {
				for i := 1; i < k; i++ {
					off++

					if b, ok := <-ch; ok {
						buf[i] = b
					} else {
						break
					}

					l++
				}
			}

			r, _ := utf8.DecodeRune(buf[:l])

			if r != utf8.RuneError && unicode.IsPrint(r) {
				rs = append(rs, r)
			} else {
				flush()
			}
		}
	}
}
