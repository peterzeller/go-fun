package list

import (
	"github.com/peterzeller/go-fun/v2/equality"
	"github.com/peterzeller/go-fun/v2/iterable"
	"github.com/peterzeller/go-fun/v2/slice"
	"github.com/peterzeller/go-fun/v2/zero"
)

// List is an immutable data type which is backed by a slice.
type List[T any] struct {
	slice []T
}

// At returns the element at the given position
func (l List[T]) At(i int) T {
	return l.slice[i]
}

// Iterator for the list.
func (l List[T]) Iterator() iterable.Iterator[T] {
	state := 0
	return iterable.Fun[T](func() (T, bool) {
		if state >= len(l.slice) {
			return zero.Value[T](), false
		}
		res := l.slice[state]
		state++
		return res, true
	})
}

// Length of the list.
func (l List[T]) Length() int {
	return len(l.slice)
}

// Create a new list
func New[T any](elems ...T) List[T] {
	s := make([]T, len(elems))
	copy(s, elems)
	return List[T]{slice: s}
}

// Append another list to this list.
func (l List[T]) Append(r List[T]) List[T] {
	if l.Length() == 0 {
		return r
	}
	if r.Length() == 0 {
		return l
	}
	s := make([]T, 0, len(l.slice)+len(r.slice))
	return List[T]{slice: append(append(s, l.slice...), r.slice...)}
}

// Contains checks whether the list contains the given element.
func (l List[T]) Contains(elem T, eq equality.Equality[T]) bool {
	return slice.ContainsEq(l.slice, elem, eq)
}

// Equal checks whether this list is equal to another list
func (l List[T]) Equal(other *List[T], eq equality.Equality[T]) bool {
	return slice.Equal(l.slice, other.slice, eq)
}

// PrefixOf checks whether this list is a prefix of another list
func (l List[T]) PrefixOf(other *List[T], eq equality.Equality[T]) bool {
	return slice.PrefixOf(l.slice, other.slice, eq)
}

// Forall checks whether all elements in the lists satisfy the given condition.
func (l List[T]) Forall(cond func(T) bool) bool {
	return slice.Exists(l.slice, cond)
}

// Exists checks whether some element in the list satisfies the given condition.
func (l List[T]) Exists(cond func(T) bool) bool {
	return slice.Exists(l.slice, cond)
}

// Skip the first n element of the list
func (l List[T]) Skip(n int) List[T] {
	if n == 0 {
		return l
	}
	if n >= l.Length() {
		return New[T]()
	}
	return New(l.slice[n:]...)
}

// Limit the length of the list and take only the first n elements.
func (l List[T]) Limit(n int) List[T] {
	if n == 0 {
		return New[T]()
	}
	if n >= l.Length() {
		return l
	}
	return New(l.slice[:n]...)
}
