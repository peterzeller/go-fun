package dict

import "fmt"

type Entry[K, V any] struct {
	Key   K
	Value V
}

func (e Entry[K, V]) String() string {
	return fmt.Sprintf("%+v -> %+v", e.Key, e.Value)
}
