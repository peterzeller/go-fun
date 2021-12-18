package mutable

import (
	"fmt"

	"github.com/peterzeller/go-fun/v2/zero"
)

type Stack[A any] struct {
	slice []A
	size  int
}

func NewStack[A any](elems ...A) *Stack[A] {
	return &Stack[A]{
		slice: elems,
		size:  len(elems),
	}
}

func (s *Stack[A]) Empty() bool {
	return s.size == 0
}

func (s *Stack[A]) Pop() A {
	if s.size <= 0 {
		panic(fmt.Errorf("popping from an empty stack"))
	}
	s.size--
	res := s.slice[s.size]
	s.slice[s.size] = zero.Value[A]()
	return res
}

func (s *Stack[A]) Push(a A) {
	if s.size < len(s.slice) {
		s.slice[s.size] = a
	} else {
		s.slice = append(s.slice, a)
	}
	s.size++
}
