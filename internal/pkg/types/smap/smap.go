package smap

import (
	"bufio"
	"bytes"
	"encoding/json"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/edsrzf/mmap-go"
)

const Tab = "    "

type action func(ch chan<- String, c *chunk)

type SMap []String

type String struct {
	Nr  int    // string nr
	Grp int    // string group
	Str string // string data
}

type chunk struct {
	min int // chunk start
	max int // chunk end
}

func Map(m *mmap.MMap) *SMap {
	s := new(SMap)

	scanner := bufio.NewScanner(bytes.NewReader(*m))

	for scanner.Scan() {
		*s = append(*s, String{
			Nr:  len(*s) + 1,
			Str: scanner.Text(),
		})
	}

	return s
}

func (s *SMap) String() string {
	var sb strings.Builder

	for _, str := range *s {
		sb.WriteString(str.Str)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (s *SMap) Render() *SMap {
	return apply(func(ch chan<- String, c *chunk) {
		for _, s := range (*s)[c.min:c.max] {
			ch <- String{s.Nr, s.Grp, indent(s.Str)}
		}
	}, len(*s))
}

func (s *SMap) Format() *SMap {
	return apply(func(ch chan<- String, c *chunk) {
		var buf bytes.Buffer

		for _, s := range (*s)[c.min:c.max] {
			buf.Reset()

			if json.Indent(&buf, []byte(s.Str), "", "  ") != nil {
				ch <- String{s.Nr, s.Grp, s.Str}
				continue
			}

			for b := range bytes.SplitSeq(buf.Bytes(), []byte("\n")) {
				ch <- String{s.Nr, s.Grp, string(b)}
			}
		}
	}, len(*s))
}

func (s *SMap) Wrap(w int) *SMap {
	return apply(func(ch chan<- String, c *chunk) {
		var i = 0
		var l string

		for _, s := range (*s)[c.min:c.max] {
			i, l = 0, indent(s.Str)

			for i < len(l)-w {
				ch <- String{s.Nr, s.Grp, l[i : i+w]}
				i += w
			}

			ch <- String{s.Nr, s.Grp, l[i:]}
		}
	}, len(*s))
}

func (s *SMap) Grep(re *regexp.Regexp) *SMap {
	return apply(func(ch chan<- String, c *chunk) {
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

func (s *SMap) CanFormat() bool {
	if len(*s) == 0 {
		return false
	}

	return json.Valid([]byte((*s)[0].Str))
}

func chunks(n int) (c []*chunk) {
	m := min(runtime.NumCPU(), n)

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
		s = append(s, str)
	}

	slices.SortStableFunc(s, func(a, b String) int {
		if a.Grp != b.Grp {
			return a.Grp - b.Grp
		} else {
			return a.Nr - b.Nr
		}
	})

	return &s
}

func indent(s string) string {
	return strings.ReplaceAll(s, "\t", Tab)
}
