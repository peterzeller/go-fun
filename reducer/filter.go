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

func Distinct[A comparable, B any](r Reducer[A, B]) Reducer[A, B] {
	return func() ReducerInstance[A, B] {
		next := r()
		seen := make(map[A]bool)
		return ReducerInstance[A, B]{
			Complete: func() B {
				return next.Complete()
			},
			Step: func(a A) bool {
				if seen[a] {
					return true
				}
				seen[a] = true
				return next.Step(a)
			},
		}
	}
}
