package heap

import (
	"math"

	"github.com/hiforensics/fox/internal/pkg/text"
)

// https://gist.github.com/n2p5/4eda328b080c9f09eff928ad47228ab1
func (h *Heap) Entropy() float64 {
	var a [256]float64
	var v float64

	ch := make(chan byte, 1024)

	go h.stream(ch)

	for b := range ch {
		a[b]++
	}

	l := float64(h.Len())

	for i := range 256 {
		if a[i] != 0 {
			f := a[i] / l
			v -= f * math.Log2(f)
		}
	}

	return v / 8
}

func (h *Heap) Strings(n, m int) <-chan text.String {
	ch := make(chan byte, 1024)
	str := make(chan text.String)

	go h.stream(ch)
	go text.Carve(ch, str, n, m)

	return str
}

func (h *Heap) stream(ch chan<- byte) {
	h.RLock()

	for _, b := range *h.mmap {
		ch <- b
	}

	h.RUnlock()

	close(ch)
}
