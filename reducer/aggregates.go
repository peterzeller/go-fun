package reducer

// Reduce is a generic reduction function that processes the input from left to right.
// Initially, the state is equal to start.
// Then, combine is called for each input together with the current state.
// Finally, the state is returned.
func Reduce[A, B any](start B, combine func(B, A) B) Reducer[A, B] {
	return func() ReducerInstance[A, B] {
		state := start
		return ReducerInstance[A, B]{
			Complete: func() B {
				return state
			},
			Step: func(a A) bool {
				state = combine(state, a)
				return true
			},
		}
	}
}

// Reduce0 is like Reduce, but uses the first element in the input as the starting value.
// Returns the zero value when the input is empty.
func Reduce0[A any](combine func(A, A) A) Reducer[A, A] {
	return func() ReducerInstance[A, A] {
		first := true
		var state A
		return ReducerInstance[A, A]{
			Complete: func() A {
				return state
			},
			Step: func(a A) bool {
				if first {
					state = a
					first = false
				} else {
					state = combine(state, a)
				}
				return true
			},
		}
	}
}

// Number type for built-in numbers
type Number interface {
	byte | int | int32 | int64 | float32 | float64
}

// Sum the input numbers.
func Sum[N Number]() Reducer[N, N] {
	return Reduce(0, func(x N, y N) N {
		return x + y
	})
}

// Count the inputs.
func Count[A any]() Reducer[A, int] {
	return Reduce(0, func(x int, y A) int {
		return x + 1
	})
}

// Product calculated the multiplication of all input values.
func Product[N Number]() Reducer[N, N] {
	return func() ReducerInstance[N, N] {
		state := N(1)
		return ReducerInstance[N, N]{
			Complete: func() N {
				return state
			},
			Step: func(a N) bool {
				if a == 0 {
					// when one number is 0, the result is 0 and we can return early
					state = 0
					return false
				}
				state = state * a
				return true
			},
		}
	}
}

// Average calculates the average value from the input
func Average[N Number]() Reducer[N, float64] {
	return func() ReducerInstance[N, float64] {
		sum := 0.0
		count := 0.0
		return ReducerInstance[N, float64]{
			Complete: func() float64 {
				if sum == 0. {
					return 0.
				}
				return sum / count
			},
			Step: func(a N) bool {
				sum += float64(a)
				count += 1
				return true
			},
		}
	}
}

// Max calculates the maximum number in the input
func Max[N Number]() Reducer[N, N] {
	return Reduce0(func(a, b N) N {
		if a > b {
			return a
		}
		return b
	})
}

// Min calculates the minimum number in the input
func Min[N Number]() Reducer[N, N] {
	return Reduce0(func(a, b N) N {
		if a < b {
			return a
		}
		return b
	})
}

// Exists checks whether there is some element in the input that satisfies the given condition.
func Exists[T any](cond func(T) bool) Reducer[T, bool] {
	return func() ReducerInstance[T, bool] {
		res := false
		return ReducerInstance[T, bool]{
			Complete: func() bool {
				return res
			},
			Step: func(a T) bool {
				if cond(a) {
					res = true
					return false
				}
				return true
			},
		}
	}
}

// Forall checks whether all elements in the input satisfy the given condition.
func Forall[T any](cond func(T) bool) Reducer[T, bool] {
	return func() ReducerInstance[T, bool] {
		res := true
		return ReducerInstance[T, bool]{
			Complete: func() bool {
				return res
			},
			Step: func(a T) bool {
				if !cond(a) {
					res = false
					return false
				}
				return true
			},
		}
	}
}

// DoErr executes the function f for all elements in the input until an error occurs.
// If an error occurs, it is returned.
func DoErr[A any](f func(A) error) Reducer[A, error] {
	return func() ReducerInstance[A, error] {
		var err error
		return ReducerInstance[A, error]{
			Complete: func() error {
				return err
			},
			Step: func(a A) bool {
				err = f(a)
				return err == nil
			},
		}
	}
}

// Do executes the function f for all elements in the input.
func Do[A any](f func(A)) Reducer[A, struct{}] {
	return func() ReducerInstance[A, struct{}] {
		return ReducerInstance[A, struct{}]{
			Complete: func() struct{} {
				return struct{}{}
			},
			Step: func(a A) bool {
				f(a)
				return true
			},
		}
	}
}
