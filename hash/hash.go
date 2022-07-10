package hash

import (
	"encoding/gob"
	"hash/fnv"

	"github.com/peterzeller/go-fun/equality"
)

// EqHash compines an equals function with a Hash function
type EqHash[T any] interface {
	equality.Equality[T]
	// Hash deterministically computes an int value fom a value.
	// Ideally, distinct values should have a high chance of distinct hash values, and computation should be fast.
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

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		uintptr
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

// Gob encoding based hash code.
// This hashes the bytes of the gob encoding of T.
// Remember that this excludes private fields, pointers, and can fail at runtime for some types.
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

// Hashable are types with a Hash function.
type Hashable interface {
	// Hash deterministically computes an int value fom a value.
	// Ideally, distinct values should have a high chance of distinct hash values, and computation should be fast.
	Hash() int64
}

// EqHashable combines equality.Equal and Hashable interfaces
type EqHashable[T any] interface {
	equality.Equal[T]
	Hashable
}

// Natural uses the hash function implemented by the types Hash function.
func Natural[T EqHashable[T]]() EqHash[T] {
	return Fun[T]{
		Eq: func(a, b T) bool {
			return a.Equal(b)
		},
		H: func(a T) int64 {
			return a.Hash()
		},
	}
}

// Map creates an EqHash instance for A, given an EqHash instance for B and a function from A to B
func Map[A, B any](base EqHash[B], f func(A) B) EqHash[A] {
	return Fun[A]{
		Eq: func(a, b A) bool {
			return base.Equal(f(a), f(b))
		},
		H: func(v A) int64 {
			return base.Hash(f(v))
		},
	}
}

type Pair[A, B any] struct {
	A A
	B B
}

// PairHash creates an EqHash instance for a pair, combining two EqHash instances
func PairHash[A, B any](a EqHash[A], b EqHash[B]) EqHash[Pair[A, B]] {
	return Fun[Pair[A, B]]{
		Eq: func(x, y Pair[A, B]) bool {
			return a.Equal(x.A, y.A) && b.Equal(x.B, y.B)
		},
		H: func(v Pair[A, B]) int64 {
			return CombineHashes(a.Hash(v.A), b.Hash(v.B))
		},
	}
}

// CombineHashes calculates a combined hash for the given hash values
func CombineHashes(hashes ...int64) int64 {
	if len(hashes) == 0 {
		return 0
	}
	result := int64(1)
	for _, element := range hashes {
		result = 31*result + element
	}
	return result
}
