package reducer_test

import (
	"testing"

	"github.com/peterzeller/go-fun/reducer"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func cmpInt(a, b int) bool {
	return a < b
}

func TestSort(t *testing.T) {
	s := []int{4, 3, 7}
	require.Equal(t, []int{3, 4, 7}, reducer.ApplySlice(s, reducer.Sorted(cmpInt, reducer.ToSlice[int]())))
}

func TestSortPartial(t *testing.T) {
	s := []int{4, 3, 7, 5, 8, 9, 10}
	count := 0
	cmp := func(a, b int) bool {
		count++
		return a < b
	}

	sorted := reducer.ApplySlice(s, reducer.Sorted(cmp, reducer.Limit(1, reducer.ToSlice[int]())))
	t.Logf("count = %d", count)
	require.Equal(t, []int{3}, sorted)
	require.True(t, count <= 9)
}

func TestSortRapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		s := rapid.SliceOf(rapid.Int()).Draw(t, "slice").([]int)
		sorted := reducer.Sorted(func(x, y int) bool { return x < y }, reducer.ToSlice[int]()).ApplySlice(s)
		require.Equal(t, len(s), len(sorted))
		for i := 0; i < len(sorted)-1; i++ {
			require.LessOrEqualf(t, sorted[i], sorted[i+1], "list should be sorted: %+v", sorted)
		}
	})
}

func TestSortBig(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		s := rapid.SliceOfN(rapid.Int(), 10000, 100000).Draw(t, "slice").([]int)
		sorted := reducer.Sorted(func(x, y int) bool { return x < y }, reducer.ToSlice[int]()).ApplySlice(s)
		require.Equal(t, len(s), len(sorted))
		for i := 0; i < len(sorted)-1; i++ {
			require.LessOrEqualf(t, sorted[i], sorted[i+1], "list should be sorted: %+v", sorted)
		}
	})
}
