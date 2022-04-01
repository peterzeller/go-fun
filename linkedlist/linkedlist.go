package linkedlist

import (
	"fmt"
	"strings"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/reducer"
	"github.com/peterzeller/go-fun/zero"
)

// LinkedList represents an immutable list.
// The empty list is represented by the nil value.
type LinkedList[T any] struct {
	head T
	tail *LinkedList[T]
}

// Iterator for the list.
func (l *LinkedList[T]) Iterator() iterable.Iterator[T] {
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
func (l *LinkedList[T]) Length() int {
	state := l
	count := 0
	for state != nil {
		count++
		state = state.tail
	}
	return count
}

// Create a new list
func New[T any](elems ...T) *LinkedList[T] {
	var res *LinkedList[T]
	for i := len(elems) - 1; i >= 0; i-- {
		res = &LinkedList[T]{
			head: elems[i],
			tail: res,
		}
	}
	return res
}

// Head is the first element in the list.
// Panics when called on the empty list.
func (l *LinkedList[T]) Head() T {
	if l == nil {
		panic(fmt.Errorf("trying to get head of empty list"))
	}
	return l.head
}

// Tail returns all but the first element of the list.
func (l *LinkedList[T]) Tail() *LinkedList[T] {
	if l == nil {
		panic(fmt.Errorf("trying to get tail of empty list"))
	}
	return l.tail
}

// Append another list to this list.
func (l *LinkedList[T]) Append(r *LinkedList[T]) *LinkedList[T] {
	var prev *LinkedList[T]
	var res *LinkedList[T]
	s := l
	for s != nil {
		node := &LinkedList[T]{
			head: s.head,
			tail: r,
		}
		if res == nil {
			res = node
		}
		if prev != nil {
			prev.tail = node
		}
		prev = node
		s = s.tail
	}
	return res
}

// Contains checks whether the list contains the given element.
func (l *LinkedList[T]) Contains(elem T, eq equality.Equality[T]) bool {
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
func (l *LinkedList[T]) Equal(other *LinkedList[T], eq equality.Equality[T]) bool {
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
func (l *LinkedList[T]) PrefixOf(other *LinkedList[T], eq equality.Equality[T]) bool {
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
func (l *LinkedList[T]) Forall(cond func(T) bool) bool {
	return reducer.Apply[T](l, reducer.Forall(cond))
}

// Exists checks whether some element in the list satisfies the given condition.
func (l *LinkedList[T]) Exists(cond func(T) bool) bool {
	return reducer.Apply[T](l, reducer.Exists(cond))
}

// Skip the first n element of the list (also named Drop in other languages)
func (l *LinkedList[T]) Skip(n int) *LinkedList[T] {
	res := l
	for i := 0; i < n; i++ {
		if res == nil {
			return nil
		}
		res = res.tail
	}
	return res
}

// Limit the length of the list and take only the first n elements (also named Take in other languages).
func (l *LinkedList[T]) Limit(n int) *LinkedList[T] {
	current := l
	var resHead *LinkedList[T]
	var resTail *LinkedList[T]
	var prev *LinkedList[T]
	for i := 0; i < n; i++ {
		if current == nil {
			return resTail
		}
		resTail = &LinkedList[T]{
			head: current.head,
		}
		if resHead == nil {
			resHead = resTail
		}
		if prev != nil {
			prev.tail = resTail
		}
		prev = resTail
		current = current.tail
	}
	return resHead
}

func (l *LinkedList[T]) String() string {
	current := l
	var s strings.Builder
	s.WriteString("[")
	first := true
	for current != nil {
		if !first {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprintf("%v", current.head))
		first = false
		current = current.tail
	}
	s.WriteString("]")
	return s.String()
}

func (l *LinkedList[T]) ToSlice() []T {
	current := l
	var res []T
	for current != nil {
		res = append(res, current.head)
		current = current.tail
	}
	return res
}
