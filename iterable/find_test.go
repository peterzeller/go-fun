package iterable_test

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	x, ok := iterable.Find(iterable.FromSlice([]int{1, 2, 3, 4, 5, 6}),
		func(t int) bool {
			return t >= 5
		})
	require.True(t, ok)
	require.Equal(t, 5, x)
}

func TestFind2(t *testing.T) {
	x, ok := iterable.Find(iterable.FromSlice([]int{1, 2, 3, 4, 5, 6}),
		func(t int) bool {
			return t >= 10
		})
	require.False(t, ok)
	require.Equal(t, 0, x)
}
