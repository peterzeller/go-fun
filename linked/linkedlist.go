package linked

import (
	"fmt"
	"strings"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/reducer"
	"github.com/peterzeller/go-fun/zero"
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
		count++
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

func Cons[T any](head T, tail *List[T]) *List[T] {
	return &List[T]{
		head: head,
		tail: tail,
	}
}

// FromIterable creates a new linked list from an iterable
func FromIterable[T any](elems iterable.Iterable[T]) *List[T] {
	var resHead, prev *List[T]
	for it := iterable.Start(elems); it.HasNext(); it.Next() {
		node := &List[T]{
			head: it.Current(),
			tail: nil,
		}
		if resHead == nil {
			resHead = node
		}
		if prev != nil {
			prev.tail = node
		}
		prev = node
	}
	return resHead
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
	var res *List[T]
	s := l
	for s != nil {
		node := &List[T]{
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
		if !eq.Equal(a.head, b.head) {
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
		if !eq.Equal(a.head, b.head) {
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

// Skip the first n element of the list (also named Drop in other languages)
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

// Limit the length of the list and take only the first n elements (also named Take in other languages).
func (l *List[T]) Limit(n int) *List[T] {
	current := l
	var resHead *List[T]
	var resTail *List[T]
	var prev *List[T]
	for i := 0; i < n; i++ {
		if current == nil {
			return resTail
		}
		resTail = &List[T]{
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

func (l *List[T]) String() string {
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

func (l *List[T]) ToSlice() []T {
	current := l
	var res []T
	for current != nil {
		res = append(res, current.head)
		current = current.tail
	}
	return res
}

func (l *List[T]) FindAndRemove(cond func(T) bool) (T, *List[T], bool) {
	var prev, resultFirst *List[T]
	c := l
	for c != nil {
		if cond(c.head) {
			if prev == nil {
				resultFirst = c.tail
			} else {
				prev.tail = c.tail
			}
			return c.head, resultFirst, true
		}
		node := &List[T]{
			head: c.head,
		}
		if resultFirst == nil {
			resultFirst = node
		}
		if prev != nil {
			prev.tail = node
		}
		prev = node
		c = c.tail
	}
	return zero.Value[T](), l, false
}

// Reversed returns a new list with the elements in reversed order.
func (l *List[T]) Reversed() *List[T] {
	var res *List[T]
	for c := l; c != nil; c = c.tail {
		res = Cons(c.head, res)
	}
	return res
}
