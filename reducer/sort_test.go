package reducer_test

import (
	"testing"

	"github.com/peterzeller/go-fun/v2/reducer"
	"github.com/stretchr/testify/require"
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
