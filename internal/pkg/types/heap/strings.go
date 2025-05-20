package heap

import (
	"unicode"
	"unicode/utf8"

	"github.com/cuhsat/fox/internal/pkg/text"
)

type String struct {
	Off int
	Str string
}

func (h *Heap) Strings(min int) <-chan String {
	ch1 := make(chan byte, 1024)
	ch2 := make(chan String)

	go h.readMMap(ch1)
	go h.carveStr(ch1, ch2, min)

	return ch2
}

func (h *Heap) readMMap(ch chan<- byte) {
	defer close(ch)

	h.RLock()

	for _, c := range *h.mmap {
		ch <- c
	}

	h.RUnlock()
}

func (h *Heap) carveStr(ch <-chan byte, s chan<- String, m int) {
	var rs []rune

	buf := make([]byte, 4)
	off := 0

	flush := func(m int) {
		if len(rs) >= m {
			s <- String{
				Off: off - (len(rs) + 1),
				Str: string(rs),
			}
		}

		rs = rs[:0]
	}

	defer close(s)
	defer flush(m)

	for b := range ch {
		buf[0] = b
		off++

		if utf8.RuneStart(b) {
			l := 1
			n := 1

			if b&0x80 == 0 {
				n = 1
			} else if b&0xE0 == 0xC0 {
				n = 2
			} else if b&0xF0 == 0xE0 {
				n = 3
			} else if b&0xF8 == 0xF0 {
				n = 4
			}

			if n > 1 && n <= utf8.UTFMax {
				for i := 1; i < n; i++ {
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
				flush(m)
			}
		} else {
			if b >= text.MinASCII && b <= text.MaxASCII {
				rs = append(rs, rune(b))
			} else {
				flush(m)
			}
		}
	}
}
