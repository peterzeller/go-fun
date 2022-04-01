package iterable

func Generate[T any](start T, next func(prev T) T) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		current := start
		first := true
		return Fun[T](func() (T, bool) {
			if first {
				first = false
				return current, true
			}
			current = next(current)
			return current, true
		})
	})
}

func GenerateState[S, T any](initialState S, next func(state S) (S, T, bool)) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		state := initialState
		return Fun[T](func() (T, bool) {
			var res T
			var ok bool
			state, res, ok = next(state)
			return res, ok
		})
	})
}
