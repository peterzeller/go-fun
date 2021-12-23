package hashdict

import (
	"fmt"

	"github.com/peterzeller/go-fun/v2/dict"
	"github.com/peterzeller/go-fun/v2/hash"
	"github.com/peterzeller/go-fun/v2/iterable"
	"github.com/peterzeller/go-fun/v2/zero"
)

// adapted from https://github.com/andrewoma/dexx/blob/master/collection/src/main/java/com/github/andrewoma/dexx/collection/internal/hashmap/CompactHashMap.java

type Dict[K, V any] struct {
	root  node[K, V]
	keyEq hash.EqHash[K]
}

func New[K, V any](eq hash.EqHash[K], entries ...dict.Entry[K, V]) Dict[K, V] {
	var root node[K, V] = empty[K, V]{}
	for _, e := range entries {
		root = root.updated0(e.Key, eq.Hash(e.Key), 0, e.Value, eq)
	}
	return Dict[K, V]{root, eq}
}

func (s Dict[K, V]) KeyEq() hash.EqHash[K] {
	return s.keyEq
}

func (d Dict[K, V]) Get(key K) (V, bool) {
	return d.root.get0(key, d.keyEq.Hash(key), 0, d.keyEq)
}

func (d Dict[K, V]) GetOrZero(key K) V {
	if r, ok := d.Get(key); ok {
		return r
	}
	return zero.Value[V]()
}

func (d Dict[K, V]) GetOr(key K, defaultValue V) V {
	if r, ok := d.Get(key); ok {
		return r
	}
	return defaultValue
}

func (d Dict[K, V]) ContainsKey(key K) bool {
	_, ok := d.Get(key)
	return ok
}

func (d Dict[K, V]) Set(key K, value V) Dict[K, V] {
	newRoot := d.root.updated0(key, d.keyEq.Hash(key), 0, value, d.keyEq)
	if newRoot == nil {
		panic(fmt.Errorf("newRoot is nil"))
	}
	return Dict[K, V]{newRoot, d.keyEq}
}

func (d Dict[K, V]) Remove(key K) Dict[K, V] {
	newRoot, changed := d.root.removed0(key, d.keyEq.Hash(key), 0, d.keyEq)
	if !changed {
		return d
	}
	if newRoot == nil {
		panic(fmt.Errorf("newRoot is nil"))
	}
	return Dict[K, V]{newRoot, d.keyEq}
}

// Iterator for the dictionary
func (d Dict[K, V]) Iterator() iterable.Iterator[dict.Entry[K, V]] {
	return d.root.iterator()
}

// Keys in the dictionary.
func (d Dict[K, V]) Keys() iterable.Iterable[K] {
	return iterable.Map(func(e dict.Entry[K, V]) K { return e.Key })(d)
}

// Values in the dictionary
func (d Dict[K, V]) Values() iterable.Iterable[V] {
	return iterable.Map(func(e dict.Entry[K, V]) V { return e.Value })(d)
}

// Number of entries in the dictionary
func (d Dict[K, V]) Size() int {
	return d.root.size()
}

// MergeAll merges the given collection of entries into this dictionary.
// The given mergeFun is called for all entries that exist in either side.
// The result returned by the mergeFun determines the new value in the merged dictionary.
// If the merge function returns nil, the value is removed for the result map.
// If an entry does not exist in one map it is given as a nil value in the merge function.
func (d Dict[K, V]) MergeAll(other iterable.Iterable[dict.Entry[K, V]], mergeFun func(K, *V, *V) *V) Dict[K, V] {
	// TODO optimize mergint with dicts
	// switch otherD := other.(type) {
	// case Dict[K, V]:
	// 	if otherD.keyEq == d.keyEq {
	// 		// special merge with other hash dictionaries using the same key:
	// 		newRoot := d.root.merge(otherD.root, mergeFun, d.keyEq)
	// 		return Dict[K, V]{newRoot, d.keyEq}
	// 	}
	// }
	res := New[K, V](d.keyEq)
	keys := New[K, struct{}](d.keyEq)
	// handle entries in other
	for it := other.Iterator(); ; {
		e, ok := it.Next()
		if !ok {
			break
		}
		dv, ok := d.Get(e.Key)
		var newV *V
		if ok {
			newV = mergeFun(e.Key, &dv, &e.Value)
		} else {
			newV = mergeFun(e.Key, nil, &e.Value)
		}
		if newV != nil {
			res = res.Set(e.Key, *newV)
		}
		keys = keys.Set(e.Key, struct{}{})
	}
	// add keys that appear in d but not in other
	for it := d.Iterator(); ; {
		e, ok := it.Next()
		if !ok {
			break
		}
		if !keys.ContainsKey(e.Key) {
			newV := mergeFun(e.Key, &e.Value, nil)
			if newV != nil {
				res = res.Set(e.Key, *newV)
			}
		}
	}
	return res
}

// Merge the given values into the dictionary.
// If an entry appears on both sides, the merge function is called to determine the new value
func (d Dict[K, V]) Merge(other iterable.Iterable[dict.Entry[K, V]], mergeFun func(K, V, V) V) Dict[K, V] {
	// TODO specialize merge with dict
	// switch otherD := other.(type) {
	// case Dict[K, V]:
	// 	if otherD.keyEq == d.keyEq {
	// 		// special merge with other hash dictionaries using the same key:
	// 		mergeFun2 := func(key K, a *V, b *V) *V {
	// 			if a == nil {
	// 				return b
	// 			}
	// 			if b == nil {
	// 				return a
	// 			}
	// 			r := mergeFun(key, *a, *b)
	// 			return &r
	// 		}
	// 		newRoot := d.root.merge(otherD.root, mergeFun2, d.keyEq)
	// 		return Dict[K, V]{newRoot, d.keyEq}
	// 	}
	// }
	res := d
	// add entries from other
	for it := other.Iterator(); ; {
		e, ok := it.Next()
		if !ok {
			break
		}
		if dv, ok := d.Get(e.Key); ok {
			newV := mergeFun(e.Key, dv, e.Value)
			d = d.Set(e.Key, newV)
		} else {
			d = d.Set(e.Key, e.Value)
		}
	}
	return res
}

// Merge the given values into the dictionary.
// If an entry appears on both sides, the value from the left side is used.
func (d Dict[K, V]) MergeLeft(other iterable.Iterable[dict.Entry[K, V]]) Dict[K, V] {
	return d.Merge(other, func(k K, v1, v2 V) V {
		return v1
	})
}

// Merge the given values into the dictionary.
// If an entry appears on both sides, the value from the left side is used.
func (d Dict[K, V]) MergeRight(other iterable.Iterable[dict.Entry[K, V]]) Dict[K, V] {
	return d.Merge(other, func(k K, v1, v2 V) V {
		return v2
	})
}
