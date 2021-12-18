package reducer

import "github.com/peterzeller/go-fun/v2/iterable"

func Map[A, B, C any](f func(A) B, r Reducer[B, C]) Reducer[A, C] {
	return func() ReducerInstance[A, C] {
		next := r()
		return ReducerInstance[A, C]{
			Complete: func() C {
				return next.Complete()
			},
			Step: func(a A) bool {
				return next.Step(f(a))
			},
		}
	}
}

func FlatMap[A, B, C any](f func(A) iterable.Iterable[B], r Reducer[B, C]) Reducer[A, C] {
	return func() ReducerInstance[A, C] {
		next := r()
		return ReducerInstance[A, C]{
			Complete: func() C {
				return next.Complete()
			},
			Step: func(a A) bool {
				it := f(a).Iterator()
				for {
					b, ok := it.Next()
					if !ok {
						return true
					}
					cont := next.Step(b)
					if !cont {
						return false
					}
				}
			},
		}
	}
}
