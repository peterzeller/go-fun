package hashset

import (
	"fmt"
	"strings"

	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/dict/hashdict"
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/reducer"
)

type Set[T any] struct {
	dict hashdict.Dict[T, struct{}]
}

// New creates a new set
func New[T any](eq hash.EqHash[T], elems ...T) Set[T] {
	entries := make([]dict.Entry[T, struct{}], len(elems))
	for i, e := range elems {
		entries[i] = dict.Entry[T, struct{}]{Key: e, Value: struct{}{}}
	}
	return Set[T]{dict: hashdict.New(eq, entries...)}
}

// EqHash returns the hash instance used in the set
func (s Set[T]) EqHash() hash.EqHash[T] {
	return s.dict.KeyEq()
}

// Contains checks if the set contains an element
func (s Set[T]) Contains(elem T) bool {
	return s.dict.ContainsKey(elem)
}

// Add elements to the set
func (s Set[T]) Add(elems ...T) Set[T] {
	d := s.dict
	for _, e := range elems {
		d = d.Set(e, struct{}{})
	}
	return Set[T]{dict: d}
}

// Remove elements from the set
func (s Set[T]) Remove(elems ...T) Set[T] {
	d := s.dict
	for _, e := range elems {
		d = d.Remove(e)
	}
	return Set[T]{dict: d}
}

// Iterator for the set
func (s Set[T]) Iterator() iterable.Iterator[T] {
	return s.dict.Keys().Iterator()

}

// Union of two sets, returning elements that return in either of the sets
func (s Set[T]) Union(other iterable.Iterable[T]) Set[T] {
	if otherS, ok := other.(Set[T]); ok {
		d := s.dict.Merge(otherS.dict, func(k T, v1, v2 struct{}) struct{} { return struct{}{} })
		return Set[T]{dict: d}
	}
	d := s.dict
	for it := other.Iterator(); ; {
		x, ok := it.Next()
		if !ok {
			break
		}
		d = d.Set(x, struct{}{})
	}
	return Set[T]{dict: d}
}

// Intersect two sets, returning only elements contained in both
func (s Set[T]) Intersect(other iterable.Iterable[T]) Set[T] {
	if otherS, ok := other.(Set[T]); ok {
		d := s.dict.MergeAll(otherS.dict, hashdict.MergeOpts[T, struct{}, struct{}, struct{}]{
			Left:  nil,
			Right: nil,
			Both: func(k T, a, b struct{}) (struct{}, bool) {
				return struct{}{}, true
			},
		})
		return Set[T]{dict: d}
	}
	res := New(s.EqHash())
	for it := other.Iterator(); ; {
		x, ok := it.Next()
		if !ok {
			break
		}
		if s.Contains(x) {
			res = res.Add(x)
		}
	}
	return res
}

// Minus removes all elements from the other set from this set
func (s Set[T]) Minus(other iterable.Iterable[T]) Set[T] {
	if otherS, ok := other.(Set[T]); ok {
		d := s.dict.MergeAll(otherS.dict, hashdict.MergeOpts[T, struct{}, struct{}, struct{}]{
			Left:  func(k T, a struct{}) (struct{}, bool) { return struct{}{}, true },
			Right: nil,
			Both: func(k T, a, b struct{}) (struct{}, bool) {
				return struct{}{}, false
			},
		})
		return Set[T]{dict: d}
	}
	res := New(s.EqHash())
	for it := s.Iterator(); ; {
		x, ok := it.Next()
		if !ok {
			break
		}
		if !reducer.Apply(other, reducer.Exists(func(y T) bool { return s.EqHash().Equal(x, y) })) {
			res = res.Add(x)
		}
	}
	return res
}

// String representation of the set
func (s Set[T]) String() string {
	var res strings.Builder
	res.WriteString("[")
	first := true
	for it := iterable.Start[T](s); it.HasNext(); it.Next() {
		if !first {
			res.WriteString(", ")
		}
		res.WriteString(fmt.Sprintf("%+v", it.Current()))
		first = false
	}
	res.WriteString("]")
	return res.String()
}
