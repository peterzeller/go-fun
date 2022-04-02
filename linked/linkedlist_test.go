package linked_test

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/linked"
	"github.com/stretchr/testify/require"
)

func Test_Limit(t *testing.T) {
	list := linked.New(1, 2, 3, 4, 5, 6, 7, 8)
	require.Equal(t, []int{1, 2, 3, 4}, list.Limit(4).ToSlice())

}

func TestAppend(t *testing.T) {
	a := linked.New(1, 2, 3)
	b := linked.New(4, 5, 6)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, a.Append(b).ToSlice())
}

func TestLength(t *testing.T) {
	list := linked.New(1, 2, 3, 4, 5, 6, 7, 8)
	require.Equal(t, 8, list.Length())
}

func TestFromIterable(t *testing.T) {
	require.Equal(t, []int{1, 2, 3}, linked.FromIterable(iterable.FromSlice([]int{1, 2, 3})).ToSlice())
}
