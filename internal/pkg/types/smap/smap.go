package smap

import (
	"bytes"
	"encoding/json"
	"regexp"
	"runtime"
	"slices"
	"sync"

	"github.com/edsrzf/mmap-go"
)

const (
	LF = 0x0a
	CR = 0x0d
)

type action func(ch chan<- String, c *chunk)

type SMap []String

type String struct {
	Nr  int
	Str string
}

type chunk struct {
	min int // chunk start
	max int // chunk end
}

func Map(m *mmap.MMap) *SMap {
	s := new(SMap)
	l := len(*m)

	var i, j int
	var a, b byte

	for ; i < l; i++ {
		a = (*m)[i]
		b = (*m)[min(i+1, l-1)]

		if a == LF || (a == CR && b != LF) {
			*s = append(*s, String{
				Nr:  len(*s) + 1,
				Str: string((*m)[j:i]),
			})

			j = i + 1
		}
	}

	*s = append(*s, String{
		Nr:  len(*s) + 1,
		Str: string((*m)[j:]),
	})

	return s
}

func (s *SMap) Indent() *SMap {
	return apply(func(ch chan<- String, c *chunk) {
		var buf bytes.Buffer
		var err error

		for _, s := range (*s)[c.min:c.max] {
			buf.Reset()

			err = json.Indent(&buf, []byte(s.Str), "", "  ")

			if err != nil {
				ch <- String{s.Nr, s.Str}
				continue
			}

			for b := range bytes.SplitSeq(buf.Bytes(), []byte("\n")) {
				ch <- String{s.Nr, string(b)}
			}
		}
	}, len(*s))
}

func (s *SMap) Wrap(w int) *SMap {
	return apply(func(ch chan<- String, c *chunk) {
		var i = 0

		for _, s := range (*s)[c.min:c.max] {
			i = 0

			for i < len(s.Str)-w {
				ch <- String{s.Nr, s.Str[i : i+w]}
				i += w
			}

			ch <- String{s.Nr, s.Str[i:]}
		}
	}, len(*s))
}

func (s *SMap) Grep(b []byte) *SMap {
	return apply(func(ch chan<- String, c *chunk) {
		re, _ := regexp.Compile(string(b))

		for _, s := range (*s)[c.min:c.max] {
			if re.MatchString(s.Str) {
				ch <- s
			}
		}
	}, len(*s))
}

func (s *SMap) Find(nr int) (int, bool) {
	if s == nil {
		return 0, false
	}

	for i, str := range *s {
		if str.Nr == nr {
			return i, true
		}
	}

	return 0, false
}

func (s *SMap) Size() (w, h int) {
	if s == nil {
		return 0, 0
	}

	for _, str := range *s {
		w = max(w, len(str.Str))
	}

	h = len(*s)

	return
}

func chunks(n int) (c []*chunk) {
	m := min(runtime.GOMAXPROCS(0), n)

	for i := range m {
		c = append(c, &chunk{
			min: i * n / m,
			max: ((i + 1) * n) / m,
		})
	}

	return
}

func apply(fn action, n int) *SMap {
	ch := make(chan String, n)

	go func() {
		var wg sync.WaitGroup

		for _, c := range chunks(n) {
			wg.Add(1)

			go func() {
				fn(ch, c)
				wg.Done()
			}()
		}

		wg.Wait()

		close(ch)
	}()

	return sort(ch)
}

func sort(ch <-chan String) *SMap {
	s := make(SMap, 0)

	for str := range ch {
		s = append(s, str) // TODO: Insert already sorted!
	}

	slices.SortStableFunc(s, func(a, b String) int {
		return a.Nr - b.Nr
	})

	return &s
}
