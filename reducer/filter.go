package reducer

func Filter[A, B any](cond func(A) bool, r Reducer[A, B]) Reducer[A, B] {
	return func() ReducerInstance[A, B] {
		next := r()
		return ReducerInstance[A, B]{
			Complete: func() B {
				return next.Complete()
			},
			Step: func(a A) bool {
				if cond(a) {
					return next.Step(a)
				}
				return true
			},
		}
	}
}

func DistinctBy[A, B any, C comparable](key func(A) C, r Reducer[A, B]) Reducer[A, B] {
	return func() ReducerInstance[A, B] {
		next := r()
		seen := make(map[C]bool)
		return ReducerInstance[A, B]{
			Complete: func() B {
				return next.Complete()
			},
			Step: func(a A) bool {
				k := key(a)
				if seen[k] {
					return true
				}
				seen[k] = true
				return next.Step(a)
			},
		}
	}
}

func Distinct[A comparable, B any](r Reducer[A, B]) Reducer[A, B] {
	return DistinctBy(func(a A) A { return a }, r)
}
