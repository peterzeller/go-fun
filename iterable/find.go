package iterable

import "github.com/peterzeller/go-fun/zero"

func Find[T any](i Iterable[T], cond func(T) bool) (T, bool) {
	it := i.Iterator()
	for {
		x, ok := it.Next()
		if !ok {
			return zero.Value[T](), false
		}
		if cond(x) {
			return x, true
		}
	}
}
