package iterable

import "testing"
import "github.com/stretchr/testify/require"

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
