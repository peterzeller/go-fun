package equality

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

func Default[T comparable]() Equality[T] {
	return Fun[T](func(a, b T) bool {
		return a == b
	})
}

func Natural[T Equal[T]]() Equality[T] {
	return Fun[T](func(a, b T) bool {
		return a.Equal(b)
	})
}

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

func SliceNatural[T Equal[T]]() Equality[[]T] {
	return Slice(Natural[T]())
}
