package smap

type stack []int

func (s *stack) Push(v int) {
	*s = append(*s, v)
}

func (s *stack) Pop() (ok bool, v int) {
	l := len(*s)

	if ok = l > 0; ok {
		*s, v = (*s)[:l-1], (*s)[l-1]
	}

	return ok, v
}
