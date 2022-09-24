package iterable

type hasLength interface {
	Length() int
}

type hasSize interface {
	Size() int
}

// Length calculates the number of elements in an iterable.
// This operation takes linear time, unless the iterable implements a Length or Size method
func Length[T any](i Iterable[T]) (size int) {
	if h, ok := i.(hasLength); ok {
		return h.Length()
	}
	if h, ok := i.(hasSize); ok {
		return h.Size()
	}
	for it := Start(i); it.HasNext(); it.Next() {
		size++
	}
	return
}
