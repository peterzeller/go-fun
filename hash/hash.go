package hash

import "github.com/peterzeller/go-fun/v2/equality"

// EqHash compines an equals function with a Hash function
type EqHash[T any] interface {
	equality.Equality[T]
	Hash(v T) int64
}

type Fun[T any] struct {
	Eq func(T, T) bool
	H  func(T) int64
}

func (f Fun[T]) Equal(a, b T) bool {
	return f.Eq(a, b)
}

func (f Fun[T]) Hash(a T) int64 {
	return f.H(a)
}

var _ EqHash[int] = Fun[int]{
	Eq: func(a, b int) bool {
		return a == b
	},
	H: func(a int) int64 {
		return int64(a)
	},
}

type Number interface {
	byte | int | int32 | int64
}

func Num[T Number]() EqHash[T] {
	return Fun[T]{
		Eq: func(a, b T) bool {
			return a == b
		},
		H: func(a T) int64 {
			return int64(a)
		},
	}
}
