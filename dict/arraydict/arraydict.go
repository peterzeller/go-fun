package arraydict

import (
	"github.com/peterzeller/go-fun/v2/dict"
	"github.com/peterzeller/go-fun/v2/equality"
	"github.com/peterzeller/go-fun/v2/iterable"
	"github.com/peterzeller/go-fun/v2/zero"
)

// ArrayDict is a dict implementation for internal use.
// Since it is used in other maps, it does not store the equality class in the struct.
// Instead, the key equality type class needs to be provided with each operation.
type ArrayDict[K, V any] struct {
	entries []dict.Entry[K, V]
}

func New[K, V any](entries ...dict.Entry[K, V]) ArrayDict[K, V] {
	es := make([]dict.Entry[K, V], 0, len(entries))
	es = append(es, entries...)
	return ArrayDict[K, V]{es}
}

func (d ArrayDict[K, V]) Get(key K, keyEq equality.Equality[K]) (V, bool) {
	for _, e := range d.entries {
		if keyEq.Equal(key, e.Key) {
			return e.Value, true
		}
	}
	return zero.Value[V](), false
}

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

func (d ArrayDict[K, V]) Size() int {
	return len(d.entries)
}

func (d ArrayDict[K, V]) First() dict.Entry[K, V] {
	if len(d.entries) == 0 {
		return zero.Value[dict.Entry[K, V]]()
	}
	return d.entries[0]
}

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
