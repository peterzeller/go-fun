package iterable

func Map[A, B any](f func(A) B) func(Iterable[A]) Iterable[B] {
	return func(base Iterable[A]) Iterable[B] {
		return &mapIterable[A, B]{base, f}
	}
}

func MapIterator[A, B any](f func(A) B) func(Iterator[A]) Iterator[B] {
	return func(base Iterator[A]) Iterator[B] {
		return &mapIterator[A, B]{base, f}
	}
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
