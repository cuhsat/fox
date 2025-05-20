package heap

import (
	"bytes"
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

	go func() {
		var b bytes.Buffer

		h.RLock()

		for _, c := range *h.mmap {
			if c >= minASCII && c <= maxASCII {
				b.WriteByte(c)
			} else {
				if b.Len() >= minString {
					ch <- b.String()
				}

				b.Reset()
			}
		}

		if b.Len() >= minString {
			ch <- b.String()
		}

		h.RUnlock()

		close(ch)
	}()

	return ch
}
