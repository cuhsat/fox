package heap

import (
	"math"
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/text"
)

func (h *Heap) Entropy(n, m float64) float64 {
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

	v /= 8

	if v < n || v > m {
		return -1 // filtered
	}

	return v
}

func (h *Heap) Strings(n, m int, i bool, re *regexp.Regexp) <-chan text.String {
	ch := make(chan byte, 1024)

	str := make(chan text.String)
	ioc := make(chan text.String)

	go h.stream(ch)
	go text.Carve(ch, str, n, m)

	if !i && re == nil {
		return str
	}

	go text.Match(str, ioc, i, re)

	return ioc
}

func (h *Heap) stream(ch chan<- byte) {
	h.RLock()

	for _, b := range *h.mmap {
		ch <- b
	}

	h.RUnlock()

	close(ch)
}
