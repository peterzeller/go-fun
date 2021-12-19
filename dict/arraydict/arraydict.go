package arraydict

import (
	"github.com/peterzeller/go-fun/v2/dict"
	"github.com/peterzeller/go-fun/v2/equality"
	"github.com/peterzeller/go-fun/v2/zero"
)

type ArrayDict[K, V any] struct {
	keyEq   equality.Equality[K]
	entries []dict.Entry[K, V]
}

func New[K, V any](eq equality.Equality[K], entries ...dict.Entry[K, V]) ArrayDict[K, V] {
	es := make([]dict.Entry[K, V], 0, len(entries))
	es = append(es, entries...)
	return ArrayDict[K, V]{eq, es}
}

func (d ArrayDict[K, V]) Get(key K) (V, bool) {
	for _, e := range d.entries {
		if d.keyEq.Equal(key, e.Key) {
			return e.Value, true
		}
	}
	return zero.Value[V](), false
}

func (d ArrayDict[K, V]) Set(key K, value V) ArrayDict[K, V] {
	newEntries := make([]dict.Entry[K, V], 0, len(d.entries)+1)
	found := false
	for _, e := range d.entries {
		if !found && d.keyEq.Equal(key, e.Key) {
			newEntries = append(newEntries, dict.Entry[K, V]{Key: key, Value: value})
			found = true
		} else {
			newEntries = append(newEntries, e)
		}
	}
	if !found {
		newEntries = append(newEntries, dict.Entry[K, V]{Key: key, Value: value})
	}
	return ArrayDict[K, V]{keyEq: d.keyEq, entries: newEntries}
}
