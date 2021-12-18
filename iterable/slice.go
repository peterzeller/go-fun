package iterable

type sliceIterable[T any] struct {
	slice []T
}

type sliceIterator[T any] struct {
	slice []T
}

func FromSlice[T any](slice []T) Iterable[T] {
	return sliceIterable[T]{slice}
}

func New[T any](slice ...T) Iterable[T] {
	return sliceIterable[T]{slice}
}

func (s sliceIterable[T]) Iterator() Iterator[T] {
	return &sliceIterator[T]{s.slice}
}

func (s *sliceIterator[T]) Next() (next T, ok bool) {
	if len(s.slice) == 0 {
		return next, false
	}
	next = s.slice[0]
	s.slice = s.slice[1:]
	return next, true
}

func ToSlice[T any](i Iterable[T]) []T {
	res := make([]T, 0)
	it := i.Iterator()
	for {
		if n, ok := it.Next(); ok {
			res = append(res, n)
			continue
		}
		return res
	}
}
