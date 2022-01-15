package slice_test

import (
	"testing"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/slice"
	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	s := []int{1, 7, 42}
	require.True(t, slice.Contains(s, 7))
	require.False(t, slice.Contains(s, 8))
}

func TestContainsEq(t *testing.T) {
	s := []int{1, 7, 42}
	require.True(t, slice.ContainsEq(s, 7, equality.Default[int]()))
	require.False(t, slice.ContainsEq(s, 8, equality.Default[int]()))
}

func TestExists(t *testing.T) {
	s := []int{1, 7, 42}
	require.True(t, slice.Exists(s, func(x int) bool { return x == 7 }))
	require.False(t, slice.Exists(s, func(x int) bool { return x == 13 }))
}

func TestForall(t *testing.T) {
	s := []int{1, 7, 42}
	require.True(t, slice.Forall(s, func(x int) bool { return x <= 42 }))
	require.False(t, slice.Forall(s, func(x int) bool { return x < 10 }))
}

func TestEqual(t *testing.T) {
	require.True(t, slice.Equal([]int{1, 2}, []int{1, 2}, equality.Default[int]()))
	require.False(t, slice.Equal([]int{1, 2}, []int{2, 1}, equality.Default[int]()))
	require.False(t, slice.Equal([]int{1, 2}, []int{1, 2, 3}, equality.Default[int]()))
}

func TestPrefixOf(t *testing.T) {
	require.True(t, slice.PrefixOf([]int{1, 2}, []int{1, 2}, equality.Default[int]()))
	require.True(t, slice.PrefixOf([]int{1, 2}, []int{1, 2, 3}, equality.Default[int]()))
	require.False(t, slice.PrefixOf([]int{1, 2}, []int{2, 1}, equality.Default[int]()))
	require.False(t, slice.PrefixOf([]int{1, 2, 3}, []int{1, 2}, equality.Default[int]()))
}
