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

func TestList_FindAndRemove1(t *testing.T) {
	x, l, ok := linked.New(1, 2, 3, 4, 5, 6).FindAndRemove(func(i int) bool {
		return i == 3
	})
	require.True(t, ok)
	require.Equal(t, x, 3)
	require.Equal(t, []int{1, 2, 4, 5, 6}, l.ToSlice())
}

func TestList_FindAndRemove2(t *testing.T) {
	x, l, ok := linked.New(1, 2, 3, 4, 5, 6).FindAndRemove(func(i int) bool {
		return i == 1
	})
	require.True(t, ok)
	require.Equal(t, x, 1)
	require.Equal(t, []int{2, 3, 4, 5, 6}, l.ToSlice())
}

func TestList_FindAndRemove3(t *testing.T) {
	x, l, ok := linked.New(1, 2, 3, 4, 5, 6).FindAndRemove(func(i int) bool {
		return i == 6
	})
	require.True(t, ok)
	require.Equal(t, x, 6)
	require.Equal(t, []int{1, 2, 3, 4, 5}, l.ToSlice())
}

func TestList_FindAndRemove4(t *testing.T) {
	_, _, ok := linked.New(1, 2, 3, 4, 5, 6).FindAndRemove(func(i int) bool {
		return i == 10
	})
	require.False(t, ok)
}

func TestList_FindAndRemove5(t *testing.T) {
	_, _, ok := linked.New[int]().FindAndRemove(func(i int) bool {
		return i == 10
	})
	require.False(t, ok)
}
