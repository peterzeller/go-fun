package dict

import (
	"github.com/peterzeller/go-fun/v2/equality"
	"github.com/peterzeller/go-fun/v2/hash"
	"github.com/peterzeller/go-fun/v2/zero"
)

// adapted from https://github.com/andrewoma/dexx/blob/master/collection/src/main/java/com/github/andrewoma/dexx/collection/internal/hashmap/CompactHashMap.java

type Dict[K, V any] struct {
	root  node[K, V]
	keyEq hash.EqHash[K]
}

type Entry[K, V any] struct {
	Key   K
	Value V
}

type node[K, V any] interface {
	size() int
	get0(key K, hash int64, level int, eq equality.Equality[K]) (V, bool)
	updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V]
}

// empty node
type empty[K, V any] struct{}

// node with one element
type singleton[K, V any] struct {
	hash  int64
	entry Entry[K, V]
}

// bucket containing multiple entries that hash to the same value
type bucket[K, V any] struct {
	hash    int64
	entries []Entry[K, V]
}

type trie[K, V any] struct {
	children sparseArray[node[K, V]]
	// count all entries
	count int
}

var _ node[int, string] = empty[int, string]{}
var _ node[int, string] = singleton[int, string]{}
var _ node[int, string] = bucket[int, string]{}
var _ node[int, string] = trie[int, string]{}

func (e empty[K, V]) size() int {
	return 0
}

func (e singleton[K, V]) size() int {
	return 1
}

func (e bucket[K, V]) size() int {
	return len(e.entries)
}

func (e trie[K, V]) size() int {
	return e.count
}

func (e empty[K, V]) get0(key K, hash int64, level int, eq equality.Equality[K]) (V, bool) {
	return zero.Value[V](), false
}

func (e singleton[K, V]) get0(key K, hash int64, level int, eq equality.Equality[K]) (V, bool) {
	if e.hash == hash && eq.Equal(e.entry.Key, key) {
		return e.entry.Value, true
	}
	return zero.Value[V](), false
}

func (e bucket[K, V]) get0(key K, hash int64, level int, eq equality.Equality[K]) (V, bool) {
	if e.hash != hash {
		return zero.Value[V](), false
	}
	for _, e := range e.entries {
		if eq.Equal(e.Key, key) {
			return e.Value, true
		}
	}
	return zero.Value[V](), false
}

func (e trie[K, V]) get0(key K, hash int64, level int, eq equality.Equality[K]) (V, bool) {
	if n, ok := e.children.get(index(hash, level)); ok {
		return n.get0(key, hash, level+5, eq)
	}
	return zero.Value[V](), false
}

func (e empty[K, V]) updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V] {
	return singleton[K, V]{
		hash:  hash,
		entry: Entry[K, V]{key, value},
	}
}

func (e singleton[K, V]) updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V] {
	if hash == e.hash {
		if eq.Equal(key, e.entry.Key) {
			// replace
			return singleton[K, V]{
				hash:  hash,
				entry: Entry[K, V]{key, value},
			}
		} else {
			// hash collision -> create bucket
			return bucket[K, V]{
				hash: hash,
				entries: []Entry[K, V]{
					e.entry,
					{key, value},
				},
			}
		}
	} else {
		// hashes are different, but collision at current level -> create a deeper trie
		e2 := singleton[K, V]{
			hash:  hash,
			entry: Entry[K, V]{key, value},
		}
		return makeTrie[K, V](e.hash, e, e2.hash, e2, level, 2, eq)
	}
}

// index in the array at the given level.
// Returns a value in {0, ..., 31}
func index(hash int64, level int) int {
	return int((uint64(hash) >> level) & 0x1f)
}

func (e bucket[K, V]) updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V] {
	if hash == e.hash {
		// add to existing bucket
		entries := make([]Entry[K, V], 0, len(e.entries)+1)
		entries = append(entries, e.entries...)
		updated := false
		for i := range entries {
			if eq.Equal(entries[i].Key, key) {
				entries[i].Value = value
				updated = true
				break
			}
		}
		if !updated {
			entries = append(entries, Entry[K, V]{key, value})
		}
		return bucket[K, V]{
			hash:    hash,
			entries: entries,
		}
	}
	// if hashes are different, make a new try
	return makeTrie[K, V](e.hash, e, hash, singleton[K, V]{hash, Entry[K, V]{key, value}}, level, e.size()+1, eq)
}

func (e trie[K, V]) updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V] {
	i := index(hash, level)
	if n, ok := e.children.get(i); ok {
		// already have a node at this index -> update that node
		n2 := n.updated0(key, hash, level+5, value, eq)
		return trie[K, V]{
			children: e.children.set(i, n2),
		}
	}

	return nil
}

func makeTrie[K, V any](aHash int64, a node[K, V], bHash int64, b node[K, V], level int, size int, eq equality.Equality[K]) node[K, V] {
	indexA := index(aHash, level)
	indexB := index(bHash, level)
	if indexA != indexB {
		return trie[K, V]{
			children: newSparseArray(
				Entry[int, node[K, V]]{indexA, a},
				Entry[int, node[K, V]]{indexB, b}),
			count: size,
		}
	} else {
		return trie[K, V]{
			children: newSparseArray(
				Entry[int, node[K, V]]{indexA, makeTrie(aHash, a, bHash, b, level+5, size, eq)}),
			count: size,
		}
	}
}
