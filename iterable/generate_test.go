package iterable_test

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	parts := iterable.TakeWhile(
		func(x int) bool {
			return x > 0
		}, iterable.Generate[int](
			10,
			func(x int) int {
				return x / 2
			}))

	require.Equal(t, []int{10, 5, 2, 1}, iterable.ToSlice(parts))
}
