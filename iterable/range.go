package iterable

// Number type for built-in numbers
type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		uintptr
}

// Range of numbers from start (inclusive) to end (exclusive)
func Range[N Number](start N, end N) Iterable[N] {
	return RangeStep(start, end, 1)
}

// Range of numbers from start (inclusive) to end (inclusive)
func RangeI[N Number](start N, end N) Iterable[N] {
	return RangeIStep(start, end, 1)
}

// Range of numbers from start (inclusive) to end (exclusive), increasing by step between elements
func RangeStep[N Number](start N, end N, step N) Iterable[N] {
	return IterableFun[N](func() Iterator[N] {
		n := start
		return Fun[N](func() (N, bool) {
			if step >= 0 && n < end || step <= 0 && n > end {
				res := n
				n += step
				return res, true
			}
			return 0, false
		})
	})
}

// Range of numbers from start (inclusive) to end (inclusive), increasing by step between elements
func RangeIStep[N Number](start N, end N, step N) Iterable[N] {
	return IterableFun[N](func() Iterator[N] {
		n := start
		return Fun[N](func() (N, bool) {
			if step >= 0 && n <= end || step <= 0 && n >= end {
				res := n
				n += step
				return res, true
			}
			return 0, false
		})
	})
}
