package dict

import (
	"math/bits"

	"github.com/peterzeller/go-fun/v2/reducer"
	"github.com/peterzeller/go-fun/v2/zero"
)

// a sparse array with at most 32 entries
type sparseArray[T any] struct {
	// bitmap has a bit value of 1 where the array has an element
	bitmap uint32
	values []T
}

func newSparseArray[T any](values ...Entry[int, T]) (res sparseArray[T]) {
	res.values = make([]T, len(values))
	i := 0
	reducer.ApplySlice(values,
		reducer.Sorted(func(a, b Entry[int, T]) bool { return a.Key < b.Key },
			reducer.Do(func(e Entry[int, T]) {
				res.bitmap = res.bitmap | (1 << e.Key)
				res.values[i] = e.Value
				i++
			})))
	return
}

func newSparseArraySorted[T any](values ...Entry[int, T]) (res sparseArray[T]) {
	res.values = make([]T, len(values))
	for i, e := range values {
		res.bitmap = res.bitmap | (1 << e.Key)
		res.values[i] = e.Value
	}
	return
}

func (a sparseArray[T]) get(i int) (T, bool) {
	mask := uint32(1) << i
	if mask&a.bitmap == 0 {
		return zero.Value[T](), false
	}
	// count the numbers of bits in the bitmap that are smaller than mask to get the real index
	realIndex := bits.OnesCount32(uint32(a.bitmap & (mask - 1)))
	return a.values[realIndex], true
}

func (a sparseArray[T]) set(i int, value T) (res sparseArray[T]) {
	res.bitmap = a.bitmap
	mask := uint32(1) << i
	if res.bitmap&mask == 0 {
		// value does not exist yet
		res.bitmap = res.bitmap | mask
		realIndex := bits.OnesCount32(uint32(a.bitmap & (mask - 1)))
		newValues := make([]T, len(a.values)+1)
		for i := 0; i < realIndex; i++ {
			newValues[i] = a.values[i]
		}
		newValues[realIndex] = value
		for i := realIndex; i < len(a.values); i++ {
			newValues[i+1] = a.values[i]
		}
		res.values = newValues
	} else {
		// overwrite existing value
		newValues := make([]T, len(a.values))
		for i, v := range a.values {
			newValues[i] = v
		}
		realIndex := bits.OnesCount32(uint32(a.bitmap & (mask - 1)))
		newValues[realIndex] = value
		res.values = newValues
	}
	return
}
