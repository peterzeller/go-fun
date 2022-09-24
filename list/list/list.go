package list

import (
	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/slice"
	"github.com/peterzeller/go-fun/zero"
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

// New creates a new list
func New[T any](elems ...T) List[T] {
	s := make([]T, len(elems))
	copy(s, elems)
	return List[T]{slice: s}
}

// FromIterable creates a new list from an iterable
func FromIterable[T any](i iterable.Iterable[T]) List[T] {
	s := make([]T, 0, iterable.Length(i))
	for it := iterable.Start(i); it.HasNext(); it.Next() {
		s = append(s, it.Current())
	}

	return List[T]{slice: s}
}

// String implements the fmt.Stringer interface
func (l List[T]) String() string {
	return iterable.String[T](l)
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

// AppendElems to the end of list
func (l List[T]) AppendElems(elem ...T) List[T] {
	if l.Length() == 0 {
		return New(elem...)
	}
	if len(elem) == 0 {
		return l
	}
	s := make([]T, 0, len(l.slice)+len(elem))
	return List[T]{slice: append(append(s, l.slice...), elem...)}
}

// Contains checks whether the list contains the given element.
func (l List[T]) Contains(elem T, eq equality.Equality[T]) bool {
	return slice.ContainsEq(l.slice, elem, eq)
}

// Equal checks whether this list is equal to another list
func (l List[T]) Equal(other List[T], eq equality.Equality[T]) bool {
	return slice.Equal(l.slice, other.slice, eq)
}

// PrefixOf checks whether this list is a prefix of another list
func (l List[T]) PrefixOf(other List[T], eq equality.Equality[T]) bool {
	return slice.PrefixOf(l.slice, other.slice, eq)
}

// Forall checks whether all elements in the lists satisfy the given condition.
func (l List[T]) Forall(cond func(T) bool) bool {
	return slice.Forall(l.slice, cond)
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

// RemoveAt removes the element at the given index from the slice and returns the modified slice.
func (l List[T]) RemoveAt(index int) List[T] {
	return List[T]{
		slice: slice.RemoveAt(l.slice, index),
	}
}

// RemoveFirst removes the first occurrence of an element from the slice and returns the modified slice.
func (l List[T]) RemoveFirst(elem T, eq equality.Equality[T]) List[T] {
	return List[T]{
		slice: slice.RemoveFirst(l.slice, elem, eq),
	}
}

// RemoveAll removes all occurrences of the element from the slice and returns the modified slice.
func (l List[T]) RemoveAll(elem T, eq equality.Equality[T]) List[T] {
	return List[T]{
		slice: slice.RemoveAll(l.slice, elem, eq),
	}
}
