package iterable

import "github.com/peterzeller/go-fun/zero"

type Iterable[T any] interface {
	Iterator() Iterator[T]
}

type Iterator[T any] interface {
	Next() (T, bool)
}

type Fun[T any] func() (T, bool)

func (f Fun[T]) Next() (T, bool) {
	return f()
}

type IterableFun[T any] func() Iterator[T]

func (f IterableFun[T]) Iterator() Iterator[T] {
	return f()
}

type emptyIterable[T any] struct {
}

type emptyIterator[T any] struct {
}

func (e emptyIterator[T]) Next() (T, bool) {
	return zero.Value[T](), false
}

func (e emptyIterable[T]) Iterator() Iterator[T] {
	return emptyIterator[T]{}
}

func Empty[T any]() Iterable[T] {
	return emptyIterable[T]{}
}

func Singleton[T any](x T) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		first := true
		return Fun[T](func() (T, bool) {
			if first {
				first = false
				return x, true
			}
			return zero.Value[T](), false
		})
	})
}
