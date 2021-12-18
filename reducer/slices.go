package reducer

func ToSlice[A any]() Reducer[A, []A] {
	return func() ReducerInstance[A, []A] {
		res := make([]A, 0)
		return ReducerInstance[A, []A]{
			Complete: func() []A {
				return res
			},
			Step: func(a A) bool {
				res = append(res, a)
				return true
			},
		}
	}

}

func ApplySlice[A, B any](s []A, reducer Reducer[A, B]) B {
	return reducer.ApplySlice(s)
}
