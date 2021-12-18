package reducer

// Limit the number of elements taken from the input to n
func Limit[A, B any](n int, next Reducer[A, B]) Reducer[A, B] {
	return func() ReducerInstance[A, B] {
		r := next()
		count := 0
		return ReducerInstance[A, B]{
			Complete: func() B {
				return r.Complete()
			},
			Step: func(a A) bool {
				if count >= n {
					return false
				}
				count++
				return r.Step(a)
			},
		}
	}
}

// Skip the first n elements in the input
func Skip[A, B any](n int, next Reducer[A, B]) Reducer[A, B] {
	return func() ReducerInstance[A, B] {
		r := next()
		count := 0
		return ReducerInstance[A, B]{
			Complete: func() B {
				return r.Complete()
			},
			Step: func(a A) bool {
				if count < n {
					count++
					return true
				}
				return r.Step(a)
			},
		}
	}
}

// First only returns the first element in the input (or zero if the input is empty)
func First[A any]() Reducer[A, A] {
	return func() ReducerInstance[A, A] {
		var res A
		return ReducerInstance[A, A]{
			Complete: func() A {
				return res
			},
			Step: func(a A) bool {
				res = a
				return false
			},
		}
	}
}
