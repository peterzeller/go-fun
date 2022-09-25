package iterable

import "github.com/peterzeller/go-fun/zero"

// Find an element in an iterable
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
