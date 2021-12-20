package hashdict

import (
	"fmt"

	"github.com/peterzeller/go-fun/v2/dict"
	"github.com/peterzeller/go-fun/v2/hash"
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

func (d Dict[K, V]) Get(key K) (V, bool) {
	return d.root.get0(key, d.keyEq.Hash(key), 0, d.keyEq)
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
