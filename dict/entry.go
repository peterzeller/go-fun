package dict

import "fmt"

// Entry in a dictionary
type Entry[K, V any] struct {
	Key   K
	Value V
}

func (e Entry[K, V]) String() string {
	return fmt.Sprintf("%+v -> %+v", e.Key, e.Value)
}

// E is a shorthand for creating an Entry
func E[K, V any](k K, v V) Entry[K, V] {
	return Entry[K, V]{Key: k, Value: v}
}
