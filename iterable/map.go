package iterable

import "github.com/peterzeller/go-fun/zero"

func Map[A, B any](base Iterable[A], f func(A) B) Iterable[B] {
	return &mapIterable[A, B]{base, f}
}

func MapIterator[A, B any](base Iterator[A], f func(A) B) Iterator[B] {
	return &mapIterator[A, B]{base, f}
}

type mapIterable[A, B any] struct {
	base Iterable[A]
	f    func(A) B
}

type mapIterator[A, B any] struct {
	base Iterator[A]
	f    func(A) B
}

func (i *mapIterable[A, B]) Iterator() Iterator[B] {
	return &mapIterator[A, B]{i.base.Iterator(), i.f}
}

func (i *mapIterator[A, B]) Next() (B, bool) {
	if a, ok := i.base.Next(); ok {
		return i.f(a), true
	}
	var b B
	return b, false
}

func FlatMap[A, B any](base Iterable[A], f func(A) Iterable[B]) Iterable[B] {
	return IterableFun[B](func() Iterator[B] {
		it := base.Iterator()
		var current Iterator[B]
		return Fun[B](func() (B, bool) {
			for {
				if current == nil {
					a, ok := it.Next()
					if !ok {
						return zero.Value[B](), false
					}
					current = f(a).Iterator()
				}
				b, ok := current.Next()
				if ok {
					return b, true
				}
				current = nil
			}
		})
	})
}
