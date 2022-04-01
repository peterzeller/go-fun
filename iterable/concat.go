package iterable

import "github.com/peterzeller/go-fun/zero"

func Concat[T any](iterables ...Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		pos := 0
		var current Iterator[T]
		return Fun[T](func() (T, bool) {
			for {
				if current == nil {
					if pos >= len(iterables) {
						return zero.Value[T](), false
					}
					current = iterables[pos].Iterator()
					pos++
				}
				n, ok := current.Next()
				if ok {
					return n, true
				}
				current = nil
			}
		})
	})
}

func ConcatIterators[T any](iterators ...Iterator[T]) Iterator[T] {
	pos := 0
	return Fun[T](func() (T, bool) {
		for pos < len(iterators) {
			n, ok := iterators[pos].Next()
			if ok {
				return n, true
			}
			pos++
		}
		return zero.Value[T](), false
	})
}
