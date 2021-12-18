package iterable

import "testing"
import "github.com/stretchr/testify/require"

func TestWhere(t *testing.T) {

	a := FromSlice([]int{1, 2, 3, 4, 5, 6})

	isEven := func(x int) bool { return x%2 == 0 }
	require.Equal(t, []int{2, 4, 6}, ToSlice(Where(isEven)(a)))

}
