package arraydict

import (
	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
)

// ArrayDict is a dict implementation for internal use.
// Since it is used in other maps, it does not store the equality class in the struct.
// Instead, the key equality type class needs to be provided with each operation.
type ArrayDict[K, V any] struct {
	entries []dict.Entry[K, V]
}

// New creates a new dictionary.
func New[K, V any](entries ...dict.Entry[K, V]) ArrayDict[K, V] {
	es := make([]dict.Entry[K, V], 0, len(entries))
	es = append(es, entries...)
	return ArrayDict[K, V]{es}
}

// Get returns the value for a given key.
func (d ArrayDict[K, V]) Get(key K, keyEq equality.Equality[K]) (V, bool) {
	for _, e := range d.entries {
		if keyEq.Equal(key, e.Key) {
			return e.Value, true
		}
	}
	return zero.Value[V](), false
}

// ContainsKey checks whether the dictionary contains the given key
func (d ArrayDict[K, V]) ContainsKey(key K, keyEq equality.Equality[K]) bool {
	_, r := d.Get(key, keyEq)
	return r
}

// Set returns an updated version of the dictionary
func (d ArrayDict[K, V]) Set(key K, value V, keyEq equality.Equality[K]) ArrayDict[K, V] {
	newEntries := make([]dict.Entry[K, V], 0, len(d.entries)+1)
	found := false
	for _, e := range d.entries {
		if !found && keyEq.Equal(key, e.Key) {
			newEntries = append(newEntries, dict.Entry[K, V]{Key: key, Value: value})
			found = true
		} else {
			newEntries = append(newEntries, e)
		}
	}
	if !found {
		newEntries = append(newEntries, dict.Entry[K, V]{Key: key, Value: value})
	}
	return ArrayDict[K, V]{entries: newEntries}
}

// Remove returns an updated dictionaries with one entry removed.
// It returns a boolean stating whether the key was found in the original map.
func (d ArrayDict[K, V]) Remove(key K, keyEq equality.Equality[K]) (ArrayDict[K, V], bool) {
	index := -1
	for i, e := range d.entries {
		if keyEq.Equal(key, e.Key) {
			index = i
			break
		}
	}
	if index == -1 {
		// not found -> unchanged
		return d, false
	}
	newEntries := make([]dict.Entry[K, V], 0, len(d.entries)-1)
	newEntries = append(append(newEntries, d.entries[:index]...), d.entries[index+1:]...)
	return ArrayDict[K, V]{entries: newEntries}, true
}

// Size returns the number of elements in the dictionary
func (d ArrayDict[K, V]) Size() int {
	return len(d.entries)
}

// First returns the first entry in the dictionary
func (d ArrayDict[K, V]) First() dict.Entry[K, V] {
	if len(d.entries) == 0 {
		return zero.Value[dict.Entry[K, V]]()
	}
	return d.entries[0]
}

// Iterator for the dictionary
func (d ArrayDict[K, V]) Iterator() iterable.Iterator[dict.Entry[K, V]] {
	pos := 0
	return iterable.Fun[dict.Entry[K, V]](func() (dict.Entry[K, V], bool) {
		if pos < len(d.entries) {
			res := d.entries[pos]
			pos++
			return res, true
		}
		return zero.Value[dict.Entry[K, V]](), false
	})
}

// FilterMap transforms the values in a map and filters them
func FilterMap[K, A, B any](d ArrayDict[K, A], f func(K, A) (B, bool)) ArrayDict[K, B] {
	res := make([]dict.Entry[K, B], 0)
	for _, e := range d.entries {
		newV, keep := f(e.Key, e.Value)
		if keep {
			res = append(res, dict.Entry[K, B]{Key: e.Key, Value: newV})
		}
	}
	return ArrayDict[K, B]{res}
}

// MergeLeft merges a dictionary and prefers the value from the left dictionary if both share a key.
func (d ArrayDict[K, V]) MergeLeft(right iterable.Iterable[dict.Entry[K, V]], keyEq equality.Equality[K]) ArrayDict[K, V] {
	res := d
	for it := iterable.Start(right); it.HasNext(); it.Next() {
		if !res.ContainsKey(it.Current().Key, keyEq) {
			res = res.Set(it.Current().Key, it.Current().Value, keyEq)
		}
	}
	return res
}

// MergeRight merges a dictionary and prefers the value from the left dictionary if both share a key.
func (d ArrayDict[K, V]) MergeRight(right iterable.Iterable[dict.Entry[K, V]], keyEq equality.Equality[K]) ArrayDict[K, V] {
	res := d
	for it := iterable.Start(right); it.HasNext(); it.Next() {
		res = res.Set(it.Current().Key, it.Current().Value, keyEq)
	}
	return res
}

// String representation of the dictionary
func (d ArrayDict[K, V]) String() string {
	return iterable.String[dict.Entry[K, V]](d)
}
