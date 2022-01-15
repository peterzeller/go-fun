package hash

import (
	"encoding/gob"
	"hash/fnv"

	"github.com/peterzeller/go-fun/equality"
)

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

var stringEq EqHash[string] = Fun[string]{
	Eq: func(a, b string) bool {
		return a == b
	},
	H: func(a string) int64 {
		h := fnv.New64a()
		h.Write([]byte(a))
		return int64(h.Sum64())
	},
}

func String() EqHash[string] {
	return stringEq
}

// Gob encoding based hash code
func Gob[T comparable]() EqHash[T] {
	return Fun[T]{
		Eq: func(a, b T) bool {
			return a == b
		},
		H: func(a T) int64 {
			h := fnv.New64a()
			d := gob.NewEncoder(h)
			d.Encode(a)
			return int64(h.Sum64())
		},
	}
}
