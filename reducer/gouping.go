package reducer

// GroupBy groups the input using the given key function.
// For each key, one instance of the valReducer is created to further process the values with that key.
func GroupBy[A any, K comparable, V any](key func(A) K, valReduce Reducer[A, V]) Reducer[A, map[K]V] {
	return func() ReducerInstance[A, map[K]V] {
		reducers := make(map[K]ReducerInstance[A, V])
		done := make(map[K]bool)
		return ReducerInstance[A, map[K]V]{
			Complete: func() map[K]V {
				res := make(map[K]V)
				for k, r := range reducers {
					res[k] = r.Complete()
				}
				return res
			},
			Step: func(v A) bool {
				k := key(v)
				if done[k] {
					return true
				}
				i, ok := reducers[k]
				if !ok {
					i = valReduce()
					reducers[k] = i
				}
				cont := i.Step(v)
				if !cont {
					done[k] = true
				}
				return true
			},
		}
	}
}

// GroupByCollect groups the input using the given key function.
// The resulting map contains all values with the same key as a slice under the same key.
func GroupByCollect[V any, K comparable](key func(V) K) Reducer[V, map[K][]V] {
	return GroupBy(key, ToSlice[V]())
}

// ToMap turns the input into a map using the given key and value functions to extract key and value from elements in the input.
// If keys appear multiple times, only take the first key.
func ToMap[T any, K comparable, V any](key func(T) K, value func(T) V) Reducer[T, map[K]V] {
	return GroupBy(key, Map(value, First[V]()))
}

// ToMapId turns the input into a map using the given key function to extract a key from elements in the input.
func ToMapId[V any, K comparable](key func(V) K) Reducer[V, map[K]V] {
	return GroupBy(key, First[V]())
}
