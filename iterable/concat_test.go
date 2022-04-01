package iterable_test

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {
	it := iterable.Concat(
		iterable.FromSlice([]int{1, 2, 3}),
		iterable.FromSlice([]int{}),
		iterable.FromSlice([]int{4, 5}),
	)

	require.Equal(t, []int{1, 2, 3, 4, 5}, iterable.ToSlice(it))
}

func TestConcatIterators(t *testing.T) {
	it := iterable.ConcatIterators(
		iterable.FromSlice([]int{1, 2, 3}).Iterator(),
		iterable.FromSlice([]int{}).Iterator(),
		iterable.FromSlice([]int{4, 5}).Iterator(),
	)

	require.Equal(t, []int{1, 2, 3, 4, 5}, iterable.IteratorToSlice(it))
}
