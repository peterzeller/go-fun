package iterable

import "testing"
import "github.com/stretchr/testify/require"

func TestMap(t *testing.T) {

	a := FromSlice([]int{1, 2, 3, 4, 5, 6})

	require.Equal(t, []int{1, 4, 9, 16, 25, 36}, ToSlice(Map(func(x int) int { return x * x })(a)))

}
