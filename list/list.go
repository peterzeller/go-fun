package list

import (
	"fmt"

	"github.com/peterzeller/go-fun/v2/equality"
	"github.com/peterzeller/go-fun/v2/iterable"
	"github.com/peterzeller/go-fun/v2/reducer"
	"github.com/peterzeller/go-fun/v2/zero"
)

// List represents an immutable list.
// The empty list is represented by the nil value.
type List[T any] struct {
	head T
	tail *List[T]
}

// Iterator for the list.
func (l *List[T]) Iterator() iterable.Iterator[T] {
	state := l
	return iterable.Fun[T](func() (T, bool) {
		if state == nil {
			return zero.Value[T](), false
		}
		res := state.head
		state = state.tail
		return res, true
	})
}

// Length of the list.
func (l *List[T]) Length() int {
	state := l
	count := 0
	for state != nil {
		state = state.tail
	}
	return count
}

// Create a new list
func New[T any](elems ...T) *List[T] {
	var res *List[T]
	for i := len(elems) - 1; i >= 0; i-- {
		res = &List[T]{
			head: elems[i],
			tail: res,
		}
	}
	return res
}

// Head is the first element in the list.
// Panics when called on the empty list.
func (l *List[T]) Head() T {
	if l == nil {
		panic(fmt.Errorf("trying to get head of empty list"))
	}
	return l.head
}

// Tail returns all but the first element of the list.
func (l *List[T]) Tail() *List[T] {
	if l == nil {
		panic(fmt.Errorf("trying to get tail of empty list"))
	}
	return l.tail
}

// Append another list to this list.
func (l *List[T]) Append(r *List[T]) *List[T] {
	var prev *List[T]
	s := l
	res := r
	for s != nil {
		res = &List[T]{
			head: s.head,
			tail: r,
		}
		if prev != nil {
			prev.tail = res
		}
		prev = res
	}
	return res
}

// Contains checks whether the list contains the given element.
func (l *List[T]) Contains(elem T, eq equality.Equality[T]) bool {
	it := l.Iterator()
	for {
		a, ok := it.Next()
		if !ok {
			return false
		}
		if eq.Equal(elem, a) {
			return true
		}
	}
}

// Equal checks whether this list is equal to another list
func (l *List[T]) Equal(other *List[T], eq equality.Equality[T]) bool {
	a := l
	b := other
	for {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		if !eq.Equal(l.head, other.head) {
			return false
		}
		a = a.tail
		b = b.tail
	}
}

// PrefixOf checks whether this list is a prefix of another list
func (l *List[T]) PrefixOf(other *List[T], eq equality.Equality[T]) bool {
	a := l
	b := other
	for {
		if a == nil {
			return true
		}
		if b == nil {
			return false
		}
		if !eq.Equal(l.head, other.head) {
			return false
		}
		a = a.tail
		b = b.tail
	}
}

// Forall checks whether all elements in the lists satisfy the given condition.
func (l *List[T]) Forall(cond func(T) bool) bool {
	return reducer.Apply[T](l, reducer.Forall(cond))
}

// Exists checks whether some element in the list satisfies the given condition.
func (l *List[T]) Exists(cond func(T) bool) bool {
	return reducer.Apply[T](l, reducer.Exists(cond))
}

// Skip the first n element of the list
func (l *List[T]) Skip(n int) *List[T] {
	res := l
	for i := 0; i < n; i++ {
		if res == nil {
			return nil
		}
		res = res.tail
	}
	return res
}

// Limit the length of the list and take only the first n elements.
func (l *List[T]) Limit(n int) *List[T] {
	current := l
	var res *List[T]
	var prev *List[T]
	for i := 0; i < n; i++ {
		if current == nil {
			return res
		}
		res = &List[T]{
			head: current.head,
		}
		if prev != nil {
			prev.tail = res
		}
		prev = res
		current = current.tail
	}
	return res
}
