package iterable

// LoopIterator is a different form of iterator better suited for use in for-loops:
//
//     for it := iterable.Start(x); it.HasNext(); it.Next() {
//        doSomethingWith(it.Current())
//     }
//
type LoopIterator[T any] struct {
	it      Iterator[T]
	current T
	ok      bool
}

func Start[T any](i Iterable[T]) LoopIterator[T] {
	it := i.Iterator()
	v, ok := it.Next()
	return LoopIterator[T]{
		it:      it,
		current: v,
		ok:      ok,
	}
}

func (l *LoopIterator[T]) HasNext() bool {
	return l.ok
}

func (l *LoopIterator[T]) Next() {
	l.current, l.ok = l.it.Next()
}

func (l *LoopIterator[T]) Current() T {
	return l.current
}

// Foreach runs the function f on each element of the iterable.
func Foreach[T any](i Iterable[T], f func(elem T)) {
	it := i.Iterator()
	for {
		elem, ok := it.Next()
		if !ok {
			return
		}
		f(elem)
	}
}
