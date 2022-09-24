package iterable

import (
	"github.com/peterzeller/go-fun/slice"
	"github.com/peterzeller/go-fun/zero"
)

func Map[A, B any](base Iterable[A], f func(A) B) Iterable[B] {
	return &mapIterable[A, B]{base, f}
}

func MapIterator[A, B any](base Iterator[A], f func(A) B) Iterator[B] {
	return &mapIterator[A, B]{base, f}
}

type mapIterable[A, B any] struct {
	base Iterable[A]
	f    func(A) B
}

// Length of a mapIterable is the same as the length of the base
func (i mapIterable[A, B]) Length() int {
	return Length(i.base)
}

type mapIterator[A, B any] struct {
	base Iterator[A]
	f    func(A) B
}

func (i mapIterable[A, B]) Iterator() Iterator[B] {
	return &mapIterator[A, B]{i.base.Iterator(), i.f}
}

func (i *mapIterator[A, B]) Next() (B, bool) {
	if a, ok := i.base.Next(); ok {
		return i.f(a), true
	}
	var b B
	return b, false
}

func FlatMap[A, B any](base Iterable[A], f func(A) Iterable[B]) Iterable[B] {
	return IterableFun[B](func() Iterator[B] {
		it := base.Iterator()
		var current Iterator[B]
		return Fun[B](func() (B, bool) {
			for {
				if current == nil {
					a, ok := it.Next()
					if !ok {
						return zero.Value[B](), false
					}
					current = f(a).Iterator()
				}
				b, ok := current.Next()
				if ok {
					return b, true
				}
				current = nil
			}
		})
	})
}

func FlatMapBreadthFirst[A, B any](base Iterable[A], f func(A) Iterable[B]) Iterable[B] {
	return IterableFun[B](func() Iterator[B] {
		it := base.Iterator()
		firstPass := true
		var iterators []Iterator[B]
		pos := 0
		return Fun[B](func() (B, bool) {
			for {
				if !firstPass && len(iterators) == 0 {
					return zero.Value[B](), false
				}
				if pos >= len(iterators) {
					if firstPass {
						// get next element from base iterator
						i, ok := it.Next()
						if ok {
							iterators = append(iterators, f(i).Iterator())
						} else {
							// no more element in base iterator
							firstPass = false
							pos = 0
							continue
						}
					} else {
						pos = 0
						continue
					}
				}
				r, ok := iterators[pos].Next()
				if ok {
					pos++
					return r, true
				}
				// remove iterator from iterators list and try with next position
				iterators = slice.RemoveAt(iterators, pos)
			}
		})
	})
}
