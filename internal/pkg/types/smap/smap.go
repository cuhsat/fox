package smap

import (
	"github.com/edsrzf/mmap-go"
)

const (
	Space = 2
)

const (
	LF = 0xa
	CR = 0xd
)

type SMap []String

type String struct {
	Nr    int
	Off   uint8
	Start int
	End   int
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
				Nr:    len(*s) + 1,
				Start: j,
				End:   i,
			})

			j = i + 1
		}
	}

	*s = append(*s, String{
		Nr:    len(*s) + 1,
		Start: j,
		End:   len(*m),
	})

	return s
}

func (s *SMap) Format(m *mmap.MMap) *SMap {
	r := new(SMap)

	pos := make(stack, 0)

	var dqt, d, i, j, l int
	var off uint8
	var ok bool

	for _, str := range *s {
		l = len(*r)

		// blank line
		if str.End-str.Start == 0 {
			*r = append(*r, str)
		}

		pos = pos[:0]
		dqt = 0
		off = 0

		for i = str.Start; i < str.End; i++ {
			switch (*m)[i] {
			case '{', '[':
				if ok, j = pos.Pop(); ok && j < i {
					add(r, str.Nr, j, i, off)
				}

				pos.Push(i + 1)

				// bracket line
				add(r, str.Nr, i, i+1, off)

				off += Space

			case '}', ']':
				if ok, j = pos.Pop(); ok && j < i {
					add(r, str.Nr, j, i, off)
				}

				off -= Space

				d = 1

				// append existing comma
				if i < str.End-1 && (*m)[i+1] == ',' {
					d += 1
				}

				// bracket line
				add(r, str.Nr, i, i+d, off)

				i += d - 1

				pos.Push(i + d)

			case ',':
				if dqt%2 != 0 {
					continue
				}

				if ok, j = pos.Pop(); ok {
					add(r, str.Nr, j, i+1, off)
				}

				pos.Push(i + 1)

			case '"':
				// parser look back
				if (*m)[max(i-1, 0)] != '\\' {
					dqt += 1
				}
			}
		}

		// normal line
		if len(*r) == l {
			*r = append(*r, str)
		}
	}

	return r
}

func (s *SMap) Wrap(w int) *SMap {
	r := new(SMap)

	for _, str := range *s {
		for str.End-str.Start > w {
			*r = append(*r, String{
				Nr:    str.Nr,
				Start: str.Start,
				End:   str.Start + w,
			})

			str.Start += w
		}

		*r = append(*r, String{
			Nr:    str.Nr,
			Start: str.Start,
			End:   str.End,
		})
	}

	return r
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
		w = max(w, str.End-str.Start)
	}

	h = len(*s)

	return
}

func add(s *SMap, n, i, j int, o uint8) {
	*s = append(*s, String{
		Nr: n, Off: o, Start: i, End: j,
	})
}
