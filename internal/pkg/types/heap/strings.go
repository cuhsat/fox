package heap

import (
	"unicode"
	"unicode/utf8"
)

const (
	minString = 3
)

const (
	minASCII = 0x20
	maxASCII = 0x7f
)

func (h *Heap) Strings() <-chan string {
	ch := make(chan string)

	rc := make(chan byte, 1024)

	go func() {
		for _, c := range *h.mmap {
			rc <- c
		}

		close(rc)
	}()

	go func() {
		var s []rune

		buf := make([]byte, 4)

		h.RLock()

		for b := range rc {
			buf[0] = b

			if utf8.RuneStart(b) {
				w := 1
				m := utf8ByteCount(b)

				if m > 1 && m <= utf8.UTFMax {
					for i := 1; i < m; i++ {
						b, ok := <-rc

						if !ok {
							break
						}

						buf[i] = b
						w++
					}
				}

				r, _ := utf8.DecodeRune(buf[:w])

				if r != utf8.RuneError && unicode.IsPrint(r) {
					s = append(s, r)
				} else {
					if len(s) >= minString {
						ch <- string(s)
					}

					s = s[:0]
				}
			} else {
				if b >= minASCII && b <= maxASCII {
					s = append(s, rune(b))
				} else {
					if len(buf) >= minString {
						ch <- string(s)
					}

					s = s[:0]
				}
			}
		}

		if len(buf) >= minString {
			ch <- string(buf)
		}

		h.RUnlock()

		close(ch)
	}()

	return ch
}

func utf8ByteCount(b byte) int {
	if b&0x80 == 0 {
		return 1
	} else if b&0xE0 == 0xC0 {
		return 2
	} else if b&0xF0 == 0xE0 {
		return 3
	} else if b&0xF8 == 0xF0 {
		return 4
	}

	return 1
}
