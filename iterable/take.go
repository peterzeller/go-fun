package iterable

import "github.com/peterzeller/go-fun/zero"

// Take the first n elements from the iterable
func Take[T any](n int, i Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		count := 0
		it := i.Iterator()
		return Fun[T](func() (T, bool) {
			if count >= n {
				return zero.Value[T](), false
			}
			count++
			return it.Next()
		})
	})
}

// TakeWhile takes elements from the iterable, while the elements match the condition
func TakeWhile[T any](cond func(T) bool, i Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		it := i.Iterator()
		active := true
		return Fun[T](func() (T, bool) {
			if !active {
				return zero.Value[T](), false
			}
			res, ok := it.Next()
			if !ok || !cond(res) {
				active = false
				return zero.Value[T](), false
			}
			return res, true
		})
	})
}
