package hashset

import (
	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/dict/hashdict"
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/reducer"
)

type Set[T any] struct {
	dict hashdict.Dict[T, struct{}]
}

func New[T any](eq hash.EqHash[T], elems ...T) Set[T] {
	entries := make([]dict.Entry[T, struct{}], len(elems))
	for i, e := range elems {
		entries[i] = dict.Entry[T, struct{}]{Key: e, Value: struct{}{}}
	}
	return Set[T]{dict: hashdict.New(eq, entries...)}
}

func (s Set[T]) EqHash() hash.EqHash[T] {
	return s.dict.KeyEq()
}

func (s Set[T]) Contains(elem T) bool {
	return s.dict.ContainsKey(elem)
}

func (s Set[T]) Add(elems ...T) Set[T] {
	d := s.dict
	for _, e := range elems {
		d = d.Set(e, struct{}{})
	}
	return Set[T]{dict: d}
}

func (s Set[T]) Remove(elems ...T) Set[T] {
	d := s.dict
	for _, e := range elems {
		d = d.Remove(e)
	}
	return Set[T]{dict: d}
}

func (s Set[T]) Iterator() iterable.Iterator[T] {
	return s.dict.Keys().Iterator()

}

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
