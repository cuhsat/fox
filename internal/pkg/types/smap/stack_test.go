package smap

import (
	"testing"
)

func TestStackPush(t *testing.T) {
	t.Run("Stack should push a new value", func(t *testing.T) {
		s := make(stack, 0)

		s.Push(1)
		s.Push(2)
		s.Push(3)

		if len(s) != 3 {
			t.Fatal("Incorrect stack size")
		}
	})
}

func TestStackPop(t *testing.T) {
	t.Run("Stack should pop the last value", func(t *testing.T) {
		s := make(stack, 0)

		s.Push(1)
		s.Push(2)
		s.Push(3)

		ok, v := s.Pop()

		if !ok || v != 3 {
			t.Fatal("Incorrect stack value")
		}
	})
}
