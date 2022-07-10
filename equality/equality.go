package equality

import "strings"

type Equal[T any] interface {
	Equal(other T) bool
}

type Equality[T any] interface {
	Equal(a, b T) bool
}

type Fun[T any] func(a, b T) bool

func (f Fun[T]) Equal(a, b T) bool {
	return f(a, b)
}

// Default equality on comparable types.
func Default[T comparable]() Equality[T] {
	return Fun[T](func(a, b T) bool {
		return a == b
	})
}

// Natural equality on types that implement the Equal interface.
func Natural[T Equal[T]]() Equality[T] {
	return Fun[T](func(a, b T) bool {
		return a.Equal(b)
	})
}

// Slice equality, given equality for slice elements.
func Slice[T any](e Equality[T]) Equality[[]T] {
	return Fun[[]T](func(a, b []T) bool {
		if len(a) != len(b) {
			return false
		}
		for i := 0; i < len(a); i++ {
			if !e.Equal(a[i], b[i]) {
				return false
			}
		}
		return true
	})
}

// SliceNatural is equivalent to Slice(Natural[T]())
func SliceNatural[T Equal[T]]() Equality[[]T] {
	return Slice(Natural[T]())
}

// StringIgnoreCase implements equality while treating upper- and lower-case letters as equivalent (using strings.EqualFold)
func StringIgnoreCase() Equality[string] {
	return Fun[string](strings.EqualFold)
}
