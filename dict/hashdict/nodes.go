package hashdict

import (
	"github.com/peterzeller/go-fun/v2/dict"
	"github.com/peterzeller/go-fun/v2/dict/arraydict"
	"github.com/peterzeller/go-fun/v2/equality"
	"github.com/peterzeller/go-fun/v2/hash"
	"github.com/peterzeller/go-fun/v2/zero"
)

type node[K, V any] interface {
	size() int
	get0(key K, hash int64, level int, eq equality.Equality[K]) (V, bool)
	updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V]
	removed0(key K, hash int64, level int, eq hash.EqHash[K]) (node[K, V], bool)
	first() (*dict.Entry[K, V], int64)
}

// empty node
type empty[K, V any] struct{}

// node with one element
type singleton[K, V any] struct {
	hash  int64
	entry dict.Entry[K, V]
}

// bucket containing multiple entries that hash to the same value
type bucket[K, V any] struct {
	hash    int64
	entries arraydict.ArrayDict[K, V]
}

type trie[K, V any] struct {
	// sparse array of size 32
	// the index is determined by 5 bits from the hash of the key, where
	// the hash is shifted by 5*depth of the node
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
	return e.entries.Size()
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
	return e.entries.Get(key, eq)
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
		entry: dict.Entry[K, V]{Key: key, Value: value},
	}
}

func (e singleton[K, V]) updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V] {
	if hash == e.hash {
		if eq.Equal(key, e.entry.Key) {
			// replace
			return singleton[K, V]{
				hash:  hash,
				entry: dict.Entry[K, V]{Key: key, Value: value},
			}
		} else {
			// hash collision -> create bucket
			return bucket[K, V]{
				hash:    hash,
				entries: arraydict.New(e.entry, dict.Entry[K, V]{Key: key, Value: value}),
			}
		}
	} else {
		// hashes are different, but collision at current level -> create a deeper trie
		e2 := singleton[K, V]{
			hash:  hash,
			entry: dict.Entry[K, V]{Key: key, Value: value},
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
		newEntries := e.entries.Set(key, value, eq)
		return bucket[K, V]{
			hash:    hash,
			entries: newEntries,
		}
	}
	// if hashes are different, make a new try
	return makeTrie[K, V](e.hash, e, hash, singleton[K, V]{hash, dict.Entry[K, V]{Key: key, Value: value}}, level, e.size()+1, eq)
}

func (e trie[K, V]) updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V] {
	i := index(hash, level)
	if n, ok := e.children.get(i); ok {
		// already have a node at this index -> update that node
		n2 := n.updated0(key, hash, level+5, value, eq)
		return trie[K, V]{
			children: e.children.set(i, n2),
			count:    e.count + (n2.size() - n.size()),
		}
	}
	// no node at the given index yet -> add singleton entry
	return trie[K, V]{
		children: e.children.set(i, singleton[K, V]{hash, dict.Entry[K, V]{Key: key, Value: value}}),
		count:    e.count + 1,
	}
}

func makeTrie[K, V any](aHash int64, a node[K, V], bHash int64, b node[K, V], level int, size int, eq equality.Equality[K]) node[K, V] {
	indexA := index(aHash, level)
	indexB := index(bHash, level)
	if indexA != indexB {
		return trie[K, V]{
			children: newSparseArray(
				dict.Entry[int, node[K, V]]{Key: indexA, Value: a},
				dict.Entry[int, node[K, V]]{Key: indexB, Value: b}),
			count: size,
		}
	} else {
		return trie[K, V]{
			children: newSparseArray(
				dict.Entry[int, node[K, V]]{Key: indexA, Value: makeTrie(aHash, a, bHash, b, level+5, size, eq)}),
			count: size,
		}
	}
}

func (e empty[K, V]) removed0(key K, hash int64, level int, eq hash.EqHash[K]) (node[K, V], bool) {
	return e, false
}

func (e singleton[K, V]) removed0(key K, hash int64, level int, eq hash.EqHash[K]) (node[K, V], bool) {
	if e.hash == hash && eq.Equal(key, e.entry.Key) {
		return empty[K, V]{}, true
	}
	return e, false
}

func (e bucket[K, V]) removed0(key K, hash int64, level int, eq hash.EqHash[K]) (node[K, V], bool) {
	if hash != e.hash {
		return e, false
	}
	newEntries, changed := e.entries.Remove(key, eq)
	if !changed {
		return e, false
	}
	switch newEntries.Size() {
	case 0:
		return empty[K, V]{}, true
	case 1:
		return singleton[K, V]{hash: hash, entry: newEntries.First()}, true
	default:
		return bucket[K, V]{
			hash:    hash,
			entries: newEntries,
		}, true
	}
}

func (e trie[K, V]) removed0(key K, hash int64, level int, eq hash.EqHash[K]) (node[K, V], bool) {
	index := index(hash, level)
	if c, ok := e.children.get(index); ok {
		newC, changed := c.removed0(key, hash, level+5, eq)
		if !changed {
			return e, false
		}
		newChildren := e.children
		if newC.size() == 0 {
			newChildren = e.children.remove(index)
			// check if we can simplify this node even more
			switch newChildren.size() {
			case 0:
				// size 0 -> simplify to empty
				return empty[K, V]{}, true
			case 1:
				// size 1 -> simplify to singleton
				singleEntry, singleHash := getFirst(newChildren)
				return singleton[K, V]{hash: singleHash, entry: *singleEntry}, true
			}
		} else {
			newChildren = e.children.set(index, newC)
		}
		return trie[K, V]{
			children: newChildren,
			count:    e.count + (newC.size() - c.size()),
		}, true
	}
	// index not in array -> unchanged
	return e, false
}

func (e empty[K, V]) first() (*dict.Entry[K, V], int64) {
	return nil, 0
}

func (e singleton[K, V]) first() (*dict.Entry[K, V], int64) {
	return &e.entry, e.hash
}

func (e bucket[K, V]) first() (*dict.Entry[K, V], int64) {
	if e.entries.Size() == 0 {
		return nil, 0
	}
	f := e.entries.First()
	return &f, e.hash
}

func (e trie[K, V]) first() (*dict.Entry[K, V], int64) {
	return getFirst(e.children)
}

func getFirst[K, V any](ar sparseArray[node[K, V]]) (*dict.Entry[K, V], int64) {
	for _, c := range ar.values {
		f, h := c.first()
		if f != nil {
			return f, h
		}
	}
	return nil, 0
}
