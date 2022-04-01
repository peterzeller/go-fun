package linkedlist_test

import (
	"testing"

	"github.com/peterzeller/go-fun/linkedlist"
	"github.com/stretchr/testify/require"
)

func Test_Limit(t *testing.T) {
	list := linkedlist.New(1, 2, 3, 4, 5, 6, 7, 8)
	require.Equal(t, []int{1, 2, 3, 4}, list.Limit(4).ToSlice())

}

func TestAppend(t *testing.T) {
	a := linkedlist.New(1, 2, 3)
	b := linkedlist.New(4, 5, 6)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, a.Append(b).ToSlice())
}

func TestLength(t *testing.T) {
	list := linkedlist.New(1, 2, 3, 4, 5, 6, 7, 8)
	require.Equal(t, 8, list.Length())
}
