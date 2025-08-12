package heap

import (
	"unicode"
	"unicode/utf8"

	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/text"
)

type String struct {
	Off int
	Str string
}

func (h *Heap) Strings(n, m int) <-chan String {
	ch := make(chan byte, 1024)
	s := make(chan String)

	go h.read(ch)
	go h.carve(ch, s, n, m)

	return s
}

func (h *Heap) read(ch chan<- byte) {
	h.RLock()

	for _, b := range *h.mmap {
		ch <- b
	}

	h.RUnlock()

	close(ch)
}

func (h *Heap) carve(ch <-chan byte, s chan<- String, n, m int) {
	var rs []rune

	flg := flags.Get().Strings
	buf := make([]byte, 4)
	off := 0

	flush := func() {
		if len(rs) >= n && len(rs) <= m {
			s <- String{
				Off: max(off-(len(rs)+1), 0),
				Str: string(rs),
			}
		}

		rs = rs[:0]
	}

	defer close(s)
	defer flush()

	for b := range ch {
		buf[0] = b
		off++

		if flg.Ascii {
			if b >= text.MinASCII && b <= text.MaxASCII {
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
					b, ok := <-ch
					off++

					if !ok {
						break
					}

					buf[i] = b
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
