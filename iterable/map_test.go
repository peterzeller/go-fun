package iterable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	a := FromSlice([]int{1, 2, 3, 4, 5, 6})
	require.Equal(t, []int{1, 4, 9, 16, 25, 36}, ToSlice(Map(a, func(x int) int { return x * x })))
}

func TestFlatMap(t *testing.T) {
	a := FromSlice([]int{1, 2, 3})
	b := FlatMap(a, func(a int) Iterable[int] {
		return FromSlice([]int{a, -a})
	})
	require.Equal(t, []int{1, -1, 2, -2, 3, -3}, ToSlice(b))
}

func TestFlatMapBreadthFirst(t *testing.T) {
	a := FromSlice([]int{1, 2, 3})
	b := FlatMapBreadthFirst(a, func(a int) Iterable[int] {
		return FromSlice([]int{a, -a})
	})
	require.Equal(t, []int{1, 2, 3, -1, -2, -3}, ToSlice(b))
}

func TestFlatMapBreadthFirst2(t *testing.T) {
	a := FromSlice([]int{10, 20, 30})
	b := FlatMapBreadthFirst(a, func(a int) Iterable[int] {
		if a == 20 {
			return New[int]()
		}
		return Range(a, a+3)
	})
	require.Equal(t, []int{10, 30, 11, 31, 12, 32}, ToSlice(b))
}

func TestFlatMapBreadthFirst3(t *testing.T) {
	a := FromSlice([]int{10, 20, 30})
	b := FlatMapBreadthFirst(a, func(a int) Iterable[int] {
		if a == 30 {
			return New[int]()
		}
		return Range(a, a+3)
	})
	require.Equal(t, []int{10, 20, 11, 21, 12, 22}, ToSlice(b))
}

func TestFlatMapBreadthFirst4(t *testing.T) {
	a := FromSlice([]int{10, 20, 30})
	b := FlatMapBreadthFirst(a, func(a int) Iterable[int] {
		if a == 10 {
			return New[int]()
		}
		return Range(a, a+3)
	})
	require.Equal(t, []int{20, 30, 21, 31, 22, 32}, ToSlice(b))
}
