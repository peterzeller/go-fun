package iterable_test

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestFun(t *testing.T) {
	r := iterable.IterableFun[int](func() iterable.Iterator[int] {
		current := 0
		return iterable.Fun[int](func() (int, bool) {
			if current < 5 {
				current++
				return current, true
			}
			return 0, false
		})
	})

	require.Equal(t, []int{1, 2, 3, 4, 5}, iterable.ToSlice[int](r))
}

func TestLoop(t *testing.T) {
	s := []int{1, 42, 7}
	is := iterable.FromSlice(s)

	var c []int
	for it := iterable.Start(is); it.HasNext(); it.Next() {
		c = append(c, it.Current())
	}
	require.Equal(t, s, c)
}

func TestMap(t *testing.T) {
	s := iterable.FromSlice([]int{1, 2, 3})
	s2 := iterable.Map(func(x int) string { return fmt.Sprintf("x%d", x) })(s)
	require.Equal(t, []string{"x1", "x2", "x3"}, iterable.ToSlice(s2))
}

func TestMapIterator(t *testing.T) {
	s := iterable.FromSlice([]int{1, 2, 3})
	it := iterable.MapIterator(func(x int) string { return fmt.Sprintf("x%d", x) })(s.Iterator())
	a, ok := it.Next()
	require.True(t, ok)
	require.Equal(t, "x1", a)
	b, ok := it.Next()
	require.True(t, ok)
	require.Equal(t, "x2", b)
	c, ok := it.Next()
	require.True(t, ok)
	require.Equal(t, "x3", c)
	_, ok = it.Next()
	require.False(t, ok)
}

func TestToString(t *testing.T) {
	s := iterable.New(1, 2, 3)
	require.Equal(t, "[1, 2, 3]", iterable.String(s))
}

func TestWhere(t *testing.T) {
	a := iterable.FromSlice([]int{1, 2, 3, 4, 5, 6})
	isEven := func(x int) bool { return x%2 == 0 }
	require.Equal(t, []int{2, 4, 6}, iterable.ToSlice(iterable.Filter(isEven)(a)))

}

func TestRange(t *testing.T) {
	require.Equal(t, []int{1, 2, 3}, iterable.ToSlice(iterable.Range(1, 4)))
}

func TestRangeI(t *testing.T) {
	require.Equal(t, []int{1, 2, 3, 4}, iterable.ToSlice(iterable.RangeI(1, 4)))
}

func TestRangeStep(t *testing.T) {
	require.Equal(t, []int{1, 4, 7, 10}, iterable.ToSlice(iterable.RangeStep(1, 13, 3)))
}

func TestRangeIStep(t *testing.T) {
	require.Equal(t, []int{1, 4, 7, 10, 13}, iterable.ToSlice(iterable.RangeIStep(1, 13, 3)))
}

func TestRangeStepRev(t *testing.T) {
	require.Equal(t, []int{13, 10, 7, 4}, iterable.ToSlice(iterable.RangeStep(13, 1, -3)))
}

func TestRangeIStepRev(t *testing.T) {
	require.Equal(t, []int{13, 10, 7, 4, 1}, iterable.ToSlice(iterable.RangeIStep(13, 1, -3)))
}
