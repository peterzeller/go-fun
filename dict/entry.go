package dict

type Entry[K, V any] struct {
	Key   K
	Value V
}
