package hashdict

import (
	"fmt"

	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/dict/arraydict"
	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
)

type node[K, V any] interface {
	size() int
	get0(key K, hash int64, level int, eq equality.Equality[K]) (V, bool)
	updated0(key K, hash int64, level int, value V, eq equality.Equality[K]) node[K, V]
	removed0(key K, hash int64, level int, eq hash.EqHash[K]) (node[K, V], bool)
	first() (*dict.Entry[K, V], int64)
	iterator() iterable.Iterator[dict.Entry[K, V]]
	checkInvariant(level int, prefix int64, eq hash.EqHash[K]) error
	fmt.Stringer
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
		return makeTrie[K, V](e.hash, e, e2.hash, e2, level, eq)
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
	return makeTrie[K, V](e.hash, e, hash, singleton[K, V]{hash, dict.Entry[K, V]{Key: key, Value: value}}, level, eq)
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

// makeTrie creates a trie from two buckets/singletons
func makeTrie[K, V any](aHash int64, a node[K, V], bHash int64, b node[K, V], level int, eq equality.Equality[K]) node[K, V] {
	if aHash == bHash {
		panic(fmt.Errorf("makeTrie called with same hash"))
	}
	size := a.size() + b.size()
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
				dict.Entry[int, node[K, V]]{Key: indexA, Value: makeTrie(aHash, a, bHash, b, level+5, eq)}),
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
				// size 1 -> check simplify to singleton
				firstNode := getFirstNode(newChildren)
				if firstNode.size() == 1 {
					singleEntry, singleHash := firstNode.first()
					return singleton[K, V]{hash: singleHash, entry: *singleEntry}, true
				}
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

func getFirstNode[K, V any](ar sparseArray[node[K, V]]) node[K, V] {
	for _, c := range ar.values {
		return c
	}
	return nil
}

func (e empty[K, V]) iterator() iterable.Iterator[dict.Entry[K, V]] {
	return iterable.Fun[dict.Entry[K, V]](func() (dict.Entry[K, V], bool) {
		return zero.Value[dict.Entry[K, V]](), false
	})
}

func (e singleton[K, V]) iterator() iterable.Iterator[dict.Entry[K, V]] {
	init := true
	return iterable.Fun[dict.Entry[K, V]](func() (dict.Entry[K, V], bool) {
		if init {
			init = false
			return e.entry, true
		}
		return zero.Value[dict.Entry[K, V]](), false
	})
}

func (e bucket[K, V]) iterator() iterable.Iterator[dict.Entry[K, V]] {
	return e.entries.Iterator()
}

func (e trie[K, V]) iterator() iterable.Iterator[dict.Entry[K, V]] {
	pos := 0
	var childIterator iterable.Iterator[dict.Entry[K, V]] = nil
	return iterable.Fun[dict.Entry[K, V]](func() (dict.Entry[K, V], bool) {
		for {
			if childIterator != nil {
				if res, ok := childIterator.Next(); ok {
					return res, true
				}
				pos++
				childIterator = nil
			}
			if pos < e.children.size() {
				childIterator = e.children.values[pos].iterator()
			} else {
				return zero.Value[dict.Entry[K, V]](), false
			}
		}
	})
}

func filterMap[K, A, B any](dictNode node[K, A], level int, eq hash.EqHash[K], f func(K, A) (B, bool)) node[K, B] {
	if f == nil {
		return empty[K, B]{}
	}
	switch e := dictNode.(type) {
	case empty[K, A]:
		return empty[K, B]{}
	case singleton[K, A]:
		newV, keep := f(e.entry.Key, e.entry.Value)
		if !keep {
			return empty[K, B]{}
		}
		return singleton[K, B]{
			hash: e.hash,
			entry: dict.Entry[K, B]{
				Key:   e.entry.Key,
				Value: newV,
			},
		}
	case bucket[K, A]:
		newEntries := arraydict.FilterMap(e.entries, f)
		switch newEntries.Size() {
		case 0:
			return empty[K, B]{}
		case 1:
			first := newEntries.First()
			return singleton[K, B]{hash: e.hash, entry: first}
		default:
			return bucket[K, B]{
				hash:    e.hash,
				entries: newEntries,
			}
		}
	case trie[K, A]:
		newChildren := sparseArrayFilterMap(e.children, func(_ int, n node[K, A]) (node[K, B], bool) {
			// recursive call
			newN := filterMap(n, level+5, eq, f)
			return newN, newN.size() > 0
		})
		switch newChildren.size() {
		case 0:
			return empty[K, B]{}
		//case 1:
		// TODO if the count is also 1, we can collapse this
		default:
			count := 0
			for _, c := range newChildren.values {
				count += c.size()
			}
			return trie[K, B]{
				children: newChildren,
				count:    count,
			}
		}
	}
	panic(fmt.Errorf("unhandled case %+v", dictNode))
}

type mergeOpts[K, A, B, C any] struct {
	eq        hash.EqHash[K]
	mergeFun  func(K, A, B) (C, bool)
	mergeFun2 func(K, B, A) (C, bool)
	// if transform functions are nil, the respective entries will be omitted
	transformA func(K, A) (C, bool)
	transformB func(K, B) (C, bool)
}

func (o mergeOpts[K, A, B, C]) swap() mergeOpts[K, B, A, C] {
	return mergeOpts[K, B, A, C]{
		eq:         o.eq,
		mergeFun:   o.mergeFun2,
		mergeFun2:  o.mergeFun,
		transformA: o.transformB,
		transformB: o.transformA,
	}
}

func (o mergeOpts[K, A, B, C]) applyA(key K, a A) (C, bool) {
	if o.transformA == nil {
		return zero.Value[C](), false
	}
	return o.transformA(key, a)
}

func (o mergeOpts[K, A, B, C]) applyB(key K, b B) (C, bool) {
	if o.transformB == nil {
		return zero.Value[C](), false
	}
	return o.transformB(key, b)
}

func trieArrayToNode[K, V any](ar sparseArray[node[K, V]]) node[K, V] {
	if ar.size() == 0 {
		return empty[K, V]{}
	}
	count := 0
	for _, c := range ar.values {
		count += c.size()
	}
	return trie[K, V]{children: ar, count: count}
}

func hashAndDictToNode[K, V any](hash int64, d arraydict.ArrayDict[K, V]) node[K, V] {
	switch d.Size() {
	case 0:
		return empty[K, V]{}
	case 1:
		return singleton[K, V]{
			hash:  hash,
			entry: d.First(),
		}
	default:
		return bucket[K, V]{
			hash:    hash,
			entries: d,
		}
	}
}

func merge[K, A, B, C any](nodeA node[K, A], nodeB node[K, B], level int, opt mergeOpts[K, A, B, C]) node[K, C] {
	// log.Printf("merge level %d: %+v and %+v", level, nodeA, nodeB)
	// switch handles the following cases:
	//           | empty | singleton | bucket | trie
	// empty     | x     | x         | x      | x
	// singleton |       | x         | x      | x
	// bucket    |       |           | x      | x
	// trie      |       |           |        | x
	// all other cases are handled by a recursive call with swapped parameters

	switch a := nodeA.(type) {
	case empty[K, A]:
		if opt.transformB == nil {
			return empty[K, C]{}
		}
		return filterMap(nodeB, level, opt.eq, opt.transformB)
	case singleton[K, A]:
		switch b := nodeB.(type) {
		case empty[K, B]:
			return merge(nodeB, nodeA, level, opt.swap())
		case singleton[K, B]:

			if a.hash == b.hash {
				// hash collision
				if opt.eq.Equal(a.entry.Key, b.entry.Key) {
					// same key -> merge
					merged, keep := opt.mergeFun(a.entry.Key, a.entry.Value, b.entry.Value)
					if keep {
						return singleton[K, C]{hash: a.hash, entry: dict.Entry[K, C]{Key: a.entry.Key, Value: merged}}
					}
					return empty[K, C]{}
				}
			}
			aNew, keepA := opt.applyA(a.entry.Key, a.entry.Value)
			bNew, keepB := opt.applyB(b.entry.Key, b.entry.Value)
			if keepA && keepB {
				if a.hash != b.hash {
					// different hashes -> create trie
					return makeTrie[K, C](a.hash, singleton[K, C]{hash: a.hash, entry: dict.Entry[K, C]{Key: a.entry.Key, Value: aNew}},
						b.hash, singleton[K, C]{hash: b.hash, entry: dict.Entry[K, C]{Key: b.entry.Key, Value: bNew}}, level, opt.eq)
				}
				// same hashes -> create bucket
				return bucket[K, C]{
					hash: a.hash,
					entries: arraydict.New(
						dict.Entry[K, C]{Key: a.entry.Key, Value: aNew},
						dict.Entry[K, C]{Key: b.entry.Key, Value: bNew}),
				}
			}
			if keepA {
				return singleton[K, C]{hash: a.hash, entry: dict.Entry[K, C]{Key: a.entry.Key, Value: aNew}}
			}
			if keepB {
				return singleton[K, C]{hash: b.hash, entry: dict.Entry[K, C]{Key: b.entry.Key, Value: bNew}}
			}
			return empty[K, C]{}
		case bucket[K, B]:
			if a.hash == b.hash {
				// hash-collision -> update bucket
				newEntries := arraydict.FilterMap(b.entries, func(key K, bv B) (C, bool) {
					if opt.eq.Equal(key, a.entry.Key) {
						return opt.mergeFun(key, a.entry.Value, bv)
					}
					return opt.applyB(key, bv)
				})
				if opt.transformA != nil {
					if !newEntries.ContainsKey(a.entry.Key, opt.eq) {
						if newVal, keep := opt.applyA(a.entry.Key, a.entry.Value); keep {
							newEntries = newEntries.Set(a.entry.Key, newVal, opt.eq)
						}
					}
				}
				return bucket[K, C]{
					hash:    a.hash,
					entries: newEntries,
				}
			}
			// different hashes -> create trie
			aNew, keepA := opt.applyA(a.entry.Key, a.entry.Value)
			bNew := filterMap[K, B](b, level, opt.eq, opt.transformB)
			if !keepA {
				return bNew
			}
			return makeTrie[K, C](a.hash, singleton[K, C]{hash: a.hash, entry: dict.Entry[K, C]{Key: a.entry.Key, Value: aNew}},
				b.hash, bNew, level, opt.eq)
		case trie[K, B]:
			aIndex := index(a.hash, level)
			updated := false
			bNew := sparseArrayFilterMap(b.children, func(i int, n node[K, B]) (node[K, C], bool) {
				var newNode node[K, C]
				if i == aIndex {
					newNode = merge(nodeA, n, level+5, opt)
					updated = true
				} else {
					if opt.transformB == nil {
						newNode = empty[K, C]{}
					} else {
						newNode = filterMap(n, level+5, opt.eq, opt.transformB)
					}
				}
				return newNode, newNode.size() > 0
			})
			if !updated {
				if newV, ok := opt.applyA(a.entry.Key, a.entry.Value); ok {
					bNew = bNew.set(aIndex, singleton[K, C]{
						hash: a.hash,
						entry: dict.Entry[K, C]{
							Key:   a.entry.Key,
							Value: newV,
						},
					})
				}
			}
			return trieArrayToNode(bNew)
		}
	case bucket[K, A]:
		switch b := nodeB.(type) {
		case empty[K, B]:
			return merge(nodeB, nodeA, level, opt.swap())
		case singleton[K, B]:
			return merge(nodeB, nodeA, level, opt.swap())
		case bucket[K, B]:
			if a.hash == b.hash {
				// hash-collision -> update bucket
				newDict := arraydict.FilterMap(a.entries, func(key K, av A) (C, bool) {
					if bv, ok := b.entries.Get(key, opt.eq); ok {
						return opt.mergeFun(key, av, bv)
					}
					return opt.applyA(key, av)
				})
				if opt.transformB != nil {
					// add entries appearing in b but not in a
					for it := iterable.Start[dict.Entry[K, B]](b.entries); it.HasNext(); it.Next() {
						if !newDict.ContainsKey(it.Current().Key, opt.eq) {
							newV, keep := opt.transformB(it.Current().Key, it.Current().Value)
							if keep {
								newDict = newDict.Set(it.Current().Key, newV, opt.eq)
							}
						}
					}

				}
				return hashAndDictToNode(a.hash, newDict)
			}
			// different hashes -> create trie
			aNew := filterMap[K, A](a, level, opt.eq, opt.transformA)
			bNew := filterMap[K, B](b, level, opt.eq, opt.transformB)
			return makeTrie[K](a.hash, aNew, b.hash, bNew, level, opt.eq)
		case trie[K, B]:
			aIndex := index(a.hash, level)
			merged := false
			bNew := sparseArrayFilterMap(b.children, func(i int, n node[K, B]) (node[K, C], bool) {
				var newNode node[K, C]
				if i == aIndex {
					newNode = merge(nodeA, n, level+5, opt)
					merged = true
				} else {
					if opt.transformB == nil {
						newNode = empty[K, C]{}
					} else {
						newNode = filterMap(n, level+5, opt.eq, opt.transformB)
					}
				}
				return newNode, newNode.size() > 0
			})
			if !merged && opt.transformA != nil {
				bNew = bNew.set(aIndex, filterMap(nodeA, level, opt.eq, opt.transformA))
			}
			return trieArrayToNode(bNew)
		}

	case trie[K, A]:
		switch b := nodeB.(type) {
		case empty[K, B]:
			return merge(nodeB, nodeA, level, opt.swap())
		case singleton[K, B]:
			return merge(nodeB, nodeA, level, opt.swap())
		case bucket[K, B]:
			return merge(nodeB, nodeA, level, opt.swap())
		case trie[K, B]:
			newEntries := make([]dict.Entry[int, node[K, C]], 0)
			for i := 0; i < 32; i++ {
				var merged node[K, C]
				if aChild, ok := a.children.get(i); ok {
					if bChild, ok := b.children.get(i); ok {
						merged = merge(aChild, bChild, level+5, opt)
					} else {
						if opt.transformA != nil {
							merged = filterMap(aChild, level+5, opt.eq, opt.transformA)
						}
					}
				} else {
					if bChild, ok := b.children.get(i); ok {
						if opt.transformB != nil {
							merged = filterMap(bChild, level+5, opt.eq, opt.transformB)
						}
					}
				}
				if merged != nil && merged.size() > 0 {
					newEntries = append(newEntries, dict.Entry[int, node[K, C]]{Key: i, Value: merged})
				}
			}
			if len(newEntries) == 0 {
				return empty[K, C]{}
			}
			count := 0
			for _, e := range newEntries {
				count += e.Value.size()
			}
			return trie[K, C]{
				children: newSparseArraySorted(newEntries...),
				count:    count,
			}
		}
	}
	panic(fmt.Errorf("unhandled case %+v, %+v", nodeA, nodeB))
}

func (e empty[K, V]) checkInvariant(level int, prefix int64, eq hash.EqHash[K]) error {
	return nil
}

func (e singleton[K, V]) checkInvariant(level int, prefix int64, eq hash.EqHash[K]) error {
	if e.hash != eq.Hash((e.entry.Key)) {
		return fmt.Errorf("wrong hash in singleton")
	}
	if e.hash<<(64-level) != prefix<<(64-level) {
		return fmt.Errorf("prefix for hash does not match in singleton %b - (prefix %b) at level %d", uint64(e.hash), uint64(prefix), level)
	}
	return nil
}

func (e bucket[K, V]) checkInvariant(level int, prefix int64, eq hash.EqHash[K]) error {
	for it := iterable.Start[dict.Entry[K, V]](e.entries); it.HasNext(); it.Next() {
		if e.hash != eq.Hash((it.Current().Key)) {
			return fmt.Errorf("wrong hash in bucket")
		}
	}
	if e.hash<<(64-level) != prefix<<(64-level) {
		return fmt.Errorf("prefix for hash does not match in bucket %b (prefix %b) at level %d", uint64(e.hash), uint64(prefix), level)
	}
	if e.size() <= 1 {
		return fmt.Errorf("bucket with %d elements", e.size())
	}
	return nil
}

func (e trie[K, V]) checkInvariant(level int, prefix int64, eq hash.EqHash[K]) error {
	count := 0
	for it := iterable.Start[dict.Entry[int, node[K, V]]](e.children); it.HasNext(); it.Next() {
		i := it.Current().Key
		c := it.Current().Value
		count += c.size()
		err := c.checkInvariant(level+5, prefix|(int64(i)<<level), eq)
		if err != nil {
			return fmt.Errorf("invalid node at index %d: %w", i, err)
		}
	}
	if count != e.count {
		return fmt.Errorf("wrong count in trie: %d (should be %d)", e.count, count)
	}
	if count == 0 {
		return fmt.Errorf("empty trie")
	}
	return nil
}

func (e empty[K, V]) String() string {
	return "empty"
}

func (e singleton[K, V]) String() string {
	return fmt.Sprintf("singleton(%b)[%+v]", e.hash, e.entry)
}

func (e bucket[K, V]) String() string {
	return fmt.Sprintf("bucket(%b)%+v", e.hash, iterable.String[dict.Entry[K, V]](e.entries))
}

func (e trie[K, V]) String() string {
	return fmt.Sprintf("trie%+v", iterable.String[dict.Entry[int, node[K, V]]](e.children))
}
