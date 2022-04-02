package hashdict

import (
	"fmt"

	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
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

func FromMap[K comparable, V any](eq hash.EqHash[K], m map[K]V) Dict[K, V] {
	var root node[K, V] = empty[K, V]{}
	for k, v := range m {
		root = root.updated0(k, eq.Hash(k), 0, v, eq)
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
	return Dict[K, V]{newRoot, d.keyEq}
}

func (d Dict[K, V]) Remove(key K) Dict[K, V] {
	newRoot, changed := d.root.removed0(key, d.keyEq.Hash(key), 0, d.keyEq)
	if !changed {
		return d
	}
	return Dict[K, V]{newRoot, d.keyEq}
}

// Iterator for the dictionary
func (d Dict[K, V]) Iterator() iterable.Iterator[dict.Entry[K, V]] {
	return d.root.iterator()
}

// Keys in the dictionary.
func (d Dict[K, V]) Keys() iterable.Iterable[K] {
	return iterable.Map[dict.Entry[K, V], K](d, func(e dict.Entry[K, V]) K { return e.Key })
}

// Values in the dictionary
func (d Dict[K, V]) Values() iterable.Iterable[V] {
	return iterable.Map[dict.Entry[K, V], V](d, func(e dict.Entry[K, V]) V { return e.Value })
}

// Number of entries in the dictionary
func (d Dict[K, V]) Size() int {
	return d.root.size()
}

func (d Dict[K, V]) String() string {
	return iterable.String[dict.Entry[K, V]](d)
}

type MergeOpts[K, A, B, C any] struct {
	Left  func(K, A) (C, bool)
	Right func(K, B) (C, bool)
	Both  func(K, A, B) (C, bool)
}

func (o MergeOpts[K, A, B, C]) intern(eq hash.EqHash[K]) mergeOpts[K, A, B, C] {
	return mergeOpts[K, A, B, C]{
		eq:       eq,
		mergeFun: o.Both,
		mergeFun2: func(k K, b B, a A) (C, bool) {
			return o.Both(k, a, b)
		},
		transformA: o.Left,
		transformB: o.Right,
	}
}

func Merge[K, A, B, C any](left Dict[K, A], right Dict[K, B], opts MergeOpts[K, A, B, C]) Dict[K, C] {
	newRoot := merge(left.root, right.root, 0, opts.intern(left.keyEq))
	return Dict[K, C]{
		root:  newRoot,
		keyEq: left.keyEq,
	}
}

func MergeIterable[K, A, B, C any](left Dict[K, A], right iterable.Iterable[dict.Entry[K, B]], opts MergeOpts[K, A, B, C]) Dict[K, C] {
	switch rightD := right.(type) {
	case Dict[K, B]:
		// special merge with other hash dictionaries using the same key:
		// we assume here that the same equality and hash code are used
		return Merge(left, rightD, opts)
	}
	res := New[K, C](left.keyEq)
	keys := New[K, struct{}](left.keyEq)
	// handle entries in right
	for it := right.Iterator(); ; {
		e, ok := it.Next()
		if !ok {
			break
		}
		dv, ok := left.Get(e.Key)
		var newV C
		keep := false
		if ok {
			newV, keep = opts.Both(e.Key, dv, e.Value)
		} else if opts.Right != nil {
			newV, keep = opts.Right(e.Key, e.Value)
		}
		if keep {
			res = res.Set(e.Key, newV)
		}
		if opts.Left != nil {
			keys = keys.Set(e.Key, struct{}{})
		}
	}
	if opts.Left != nil {
		// add keys that appear in left but not in right
		for it := iterable.Start[dict.Entry[K, A]](left); it.HasNext(); it.Next() {
			e := it.Current()
			if !keys.ContainsKey(e.Key) {
				newV, keep := opts.Left(e.Key, e.Value)
				if keep {
					res = res.Set(e.Key, newV)
				}
			}
		}
	}
	return res
}

// MergeAll merges the given collection of entries into this dictionary.
func (d Dict[K, V]) MergeAll(other iterable.Iterable[dict.Entry[K, V]], opts MergeOpts[K, V, V, V]) Dict[K, V] {
	return MergeIterable(d, other, opts)
}

// Merge the given values into the dictionary.
// If an entry appears on both sides, the merge function is called to determine the new value
func (d Dict[K, V]) Merge(other iterable.Iterable[dict.Entry[K, V]], mergeFun func(K, V, V) V) Dict[K, V] {
	return MergeIterable(d, other, MergeOpts[K, V, V, V]{
		Left:  func(k K, a V) (V, bool) { return a, true },
		Right: func(k K, b V) (V, bool) { return b, true },
		Both:  func(k K, a V, b V) (V, bool) { return mergeFun(k, a, b), true },
	})
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

func (d Dict[K, V]) checkInvariant() error {
	if d.root == nil {
		return fmt.Errorf("root is nil")
	}
	if d.keyEq == nil {
		return fmt.Errorf("keyEq is nil")
	}
	return d.root.checkInvariant(0, 0, d.keyEq)
}

func FilterMap[K, A, B any](d Dict[K, A], f func(K, A) (B, bool)) Dict[K, B] {
	return Dict[K, B]{
		keyEq: d.keyEq,
		root:  filterMap(d.root, 0, d.keyEq, f),
	}
}

func (d Dict[K, V]) FilterMap(f func(K, V) (V, bool)) Dict[K, V] {
	return FilterMap(d, f)
}

func Map[K, A, B any](d Dict[K, A], f func(K, A) B) Dict[K, B] {
	return Dict[K, B]{
		keyEq: d.keyEq,
		root: filterMap(d.root, 0, d.keyEq, func(key K, value A) (B, bool) {
			return f(key, value), true
		}),
	}
}

func (d Dict[K, V]) Map(f func(K, V) V) Dict[K, V] {
	return Map(d, f)
}

func (d Dict[K, V]) Filter(cond func(K, V) bool) Dict[K, V] {
	return Dict[K, V]{
		keyEq: d.keyEq,
		root: filterMap(d.root, 0, d.keyEq, func(key K, value V) (V, bool) {
			return value, cond(key, value)
		}),
	}
}

var notEqual = fmt.Errorf("not equal")

func (d Dict[K, V]) Equal(other Dict[K, V], eq equality.Equality[V]) (res bool) {
	if d.Size() != other.Size() {
		return false
	}
	// we can use iterators, since the iteration order of a trie is deterministic
	// there is some optimization potential with a recursive equal function that uses reference equality of subtrees
	it1 := d.Iterator()
	it2 := other.Iterator()

	for {
		e1, ok1 := it1.Next()
		e2, ok2 := it2.Next()
		// since sizes are equal, ok1 == ok2
		if !ok1 && !ok2 {
			return true
		}
		if !d.keyEq.Equal(e1.Key, e2.Key) || !eq.Equal(e1.Value, e2.Value) {
			return false
		}
	}
}
