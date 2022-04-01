package iterable_test

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestTake(t *testing.T) {
	it := iterable.Take(2, iterable.FromSlice([]int{1, 2, 3, 4, 5}))
	require.Equal(t, []int{1, 2}, iterable.ToSlice(it))
}

func TestTakeWhile(t *testing.T) {
	it := iterable.TakeWhile(func(x int) bool {
		return x <= 3
	}, iterable.FromSlice([]int{1, 2, 3, 4, 3, 2, 1}))
	require.Equal(t, []int{1, 2, 3}, iterable.ToSlice(it))
}
