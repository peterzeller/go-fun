package iterable_test

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestTake(t *testing.T) {
	it := iterable.Take(2, iterable.FromSlice([]int{1, 2, 3, 4, 5}))
	require.Equal(t, []int{1, 2}, iterable.ToSlice(it))
}

func ExampleTakeWhile() {
	it := iterable.TakeWhile(func(x int) bool {
		return x <= 3
	}, iterable.FromSlice([]int{1, 2, 3, 4, 3, 2, 1}))
	fmt.Printf("it = %s\n", iterable.String(it))
	// output: it = [1, 2, 3]
}

func TestTakeWhile(t *testing.T) {
	tw := iterable.TakeWhile(func(x int) bool {
		return x <= 3
	}, iterable.New(1, 2, 3, 4, 3, 2, 1))
	it := tw.Iterator()
	v, ok := it.Next()
	require.Equal(t, 1, v)
	require.True(t, ok)
	v, ok = it.Next()
	require.Equal(t, 2, v)
	require.True(t, ok)
	v, ok = it.Next()
	require.Equal(t, 3, v)
	require.True(t, ok)
	v, ok = it.Next()
	require.Equal(t, 0, v)
	require.False(t, ok)
	v, ok = it.Next()
	require.Equal(t, 0, v)
	require.False(t, ok)
}
