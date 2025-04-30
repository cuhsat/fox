package smap

type stack []int

func (s *stack) Push(v int) {
	*s = append(*s, v)
}

func (s *stack) Pop() (v int) {
	l := len(*s)

	if l == 0 {
		return -1
	}

	*s, v = (*s)[:l-1], (*s)[l-1]

	return
}
