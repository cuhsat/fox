package smap

import (
	"github.com/edsrzf/mmap-go"
)

const (
	Space = 2
)

type SMap []String

type String struct {
	Nr    int
	Start int
	End   int
	Len   int
	Off   int
}

func Map(m *mmap.MMap) *SMap {
	s := new(SMap)

	var i, j int

	for ; i < len(*m); i++ {
		if (*m)[i] != '\n' {
			continue
		}

		*s = append(*s, String{
			Nr:    len(*s) + 1,
			Start: j,
			End:   i,
			Len:   i - j,
		})

		j = i + 1
	}

	*s = append(*s, String{
		Nr:    len(*s) + 1,
		Start: j,
		End:   len(*m),
		Len:   len(*m) - j,
	})

	return s
}

func (s *SMap) Format(m *mmap.MMap) *SMap {
	r := new(SMap)

	for _, str := range *s {
		l := len(*r)

		// blank line
		if str.Len == 0 {
			*r = append(*r, str)
		}

		pos := make(stack, 0)
		dqt := 0
		off := 0

		for i := str.Start; i < str.End; i++ {
			switch (*m)[i] {
			case '{', '[':
				if ok, j := pos.Pop(); ok && j < i {
					add(r, str.Nr, j, i, off)
				}

				pos.Push(i + 1)

				// bracket line
				add(r, str.Nr, i, i+1, off)

				off += Space

			case '}', ']':
				if ok, j := pos.Pop(); ok && j < i {
					add(r, str.Nr, j, i, off)
				}

				off -= Space

				d := 1

				// append existing comma
				if i < str.End-1 && (*m)[i+1] == ',' {
					d += 1
				}

				// bracket line
				add(r, str.Nr, i, i+d, off)

				i += (d - 1)

				pos.Push(i + d)

			case ',':
				if dqt%2 != 0 {
					continue
				}

				if ok, j := pos.Pop(); ok {
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
		for str.Len > w {
			*r = append(*r, String{
				Nr:    str.Nr,
				Start: str.Start,
				End:   str.Start + w,
				Len:   w,
			})

			str.Start += w
			str.Len -= w
		}

		*r = append(*r, String{
			Nr:    str.Nr,
			Start: str.Start,
			End:   str.End,
			Len:   str.Len,
		})
	}

	return r
}

func (s *SMap) Find(nr int) (bool, int) {
	for i, str := range *s {
		if str.Nr == nr {
			return true, i
		}
	}

	return false, 0
}

func (s *SMap) Size() (w, h int) {
	for _, str := range *s {
		w = max(w, str.Len)
	}

	h = len(*s)

	return
}

func add(s *SMap, n, i, j, o int) {
	*s = append(*s, String{
		Nr: n, Start: i, End: j, Len: j - i, Off: o,
	})
}
