package iterable

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
