package hashdict

import (
	"math/bits"

	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/reducer"
	"github.com/peterzeller/go-fun/zero"
)

// a sparse array with at most 32 entries
type sparseArray[T any] struct {
	// bitmap has a bit value of 1 where the array has an element
	bitmap uint32
	values []T
}

func newSparseArray[T any](values ...dict.Entry[int, T]) (res sparseArray[T]) {
	res.values = make([]T, len(values))
	i := 0
	reducer.ApplySlice(values,
		reducer.Sorted(func(a, b dict.Entry[int, T]) bool { return a.Key < b.Key },
			reducer.Do(func(e dict.Entry[int, T]) {
				res.bitmap = res.bitmap | (1 << e.Key)
				res.values[i] = e.Value
				i++
			})))
	return
}

func newSparseArraySorted[T any](values ...dict.Entry[int, T]) (res sparseArray[T]) {
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

func (a sparseArray[T]) remove(i int) sparseArray[T] {
	mask := uint32(1) << i
	if a.bitmap&mask == 0 {
		// removed index does not exist -> unchanged
		return a
	}
	var res sparseArray[T]
	// remove from bitmap
	res.bitmap = a.bitmap & (^mask)
	// remove from array
	realIndex := bits.OnesCount32(uint32(a.bitmap & (mask - 1)))
	res.values = append(append(res.values, a.values[:realIndex]...), a.values[realIndex+1:]...)
	return res
}

func (a sparseArray[T]) size() int {
	return len(a.values)
}

func sparseArrayFilterMap[A, B any](a sparseArray[A], f func(int, A) (B, bool)) sparseArray[B] {
	bitmap := uint32(0)
	values := make([]B, 0)
	j := 0
	for i := 0; i < 32; i++ {
		mask := uint32(1) << i
		if a.bitmap&mask == 0 {
			continue
		}
		newV, keep := f(i, a.values[j])
		j++
		if keep {
			bitmap = bitmap | mask
			values = append(values, newV)
		}
	}
	return sparseArray[B]{
		bitmap: bitmap,
		values: values,
	}
}

func (a sparseArray[T]) Iterator() iterable.Iterator[dict.Entry[int, T]] {
	i := 0
	return iterable.Fun[dict.Entry[int, T]](func() (dict.Entry[int, T], bool) {
		for i < 32 {
			if v, ok := a.get(i); ok {
				res := dict.Entry[int, T]{Key: i, Value: v}
				i++
				return res, true
			}
			i++
		}
		return zero.Value[dict.Entry[int, T]](), false
	})
}
