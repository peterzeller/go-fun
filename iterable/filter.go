package iterable

import "github.com/peterzeller/go-fun/zero"

func Filter[A any](base Iterable[A], cond func(A) bool) Iterable[A] {
	return &whereIterable[A]{base, cond}
}

type whereIterable[A any] struct {
	base Iterable[A]
	cond func(A) bool
}

type whereIterator[A any] struct {
	base Iterator[A]
	cond func(A) bool
}

func (i *whereIterable[A]) Iterator() Iterator[A] {
	return &whereIterator[A]{i.base.Iterator(), i.cond}
}

func (i *whereIterator[A]) Next() (A, bool) {
	for {
		if a, ok := i.base.Next(); ok {
			if i.cond(a) {
				return a, true
			}
			continue
		}
		return zero.Value[A](), false
	}
}
