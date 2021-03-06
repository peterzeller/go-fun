package reducer_test

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/reducer"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	square := reducer.Map(func(x int) int { return x * x },
		reducer.ToSlice[int]())
	require.Equal(t, []int{1, 4, 9, 16, 25, 36}, reducer.ApplySlice(s, square))
}

func TestFlatMap(t *testing.T) {
	s := []int{1, 2, 3}
	plusMinus := reducer.FlatMap(func(x int) iterable.Iterable[int] { return iterable.New(x, -x) },
		reducer.ToSlice[int]())
	require.Equal(t, []int{1, -1, 2, -2, 3, -3}, reducer.ApplySlice(s, plusMinus))
}

func TestFlatMapLimit(t *testing.T) {
	s := reducer.Apply(iterable.New(1, 2, 3),
		reducer.FlatMap(func(i int) iterable.Iterable[int] { return iterable.New(-i, i) },
			reducer.Limit(3,
				reducer.ToSlice[int]())))
	require.Equal(t, []int{-1, 1, -2}, s)
}

func TestFilter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	onlyEven := reducer.Filter(func(x int) bool { return x%2 == 0 },
		reducer.ToSlice[int]())
	require.Equal(t, []int{2, 4, 6}, reducer.ApplySlice(s, onlyEven))
}

func TestSum(t *testing.T) {
	s := []int{1, 2, 3}
	require.Equal(t, 6, reducer.ApplySlice(s, reducer.Sum[int]()))
}

func TestProd(t *testing.T) {
	s := []int{2, 3, 4}
	require.Equal(t, 24, reducer.ApplySlice(s, reducer.Product[int]()))
}

func TestProdWithZero(t *testing.T) {
	s := []int{2, 3, 0, 4, 5}
	require.Equal(t, 0, reducer.ApplySlice(s, reducer.Product[int]()))
}

func TestSkip(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	require.Equal(t, []int{3, 4, 5, 6}, reducer.ApplySlice(s, reducer.Skip(2, reducer.ToSlice[int]())))
}

func TestLimit(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	require.Equal(t, []int{1, 2, 3}, reducer.ApplySlice(s, reducer.Limit(3, reducer.ToSlice[int]())))
}

func TestSkipLimit(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	require.Equal(t, []int{3, 4, 5}, reducer.ApplySlice(s, reducer.Skip(2, reducer.Limit(3, reducer.ToSlice[int]()))))
}

func TestDistinct(t *testing.T) {
	s := []int{1, 2, 3, 2, 2, 3, 4, 1, 5}
	require.Equal(t, []int{1, 2, 3, 4, 5}, reducer.ApplySlice(s, reducer.Distinct(reducer.ToSlice[int]())))
}

type Book struct {
	Author string
	Title  string
	Year   int
}

func TestBook(t *testing.T) {
	books := []Book{
		{"A", "Q", 1990},
		{"B", "R", 2005},
		{"A", "S", 2001},
		{"B", "T", 1999},
		{"B", "U", 2021},
	}

	m := reducer.ApplySlice(books,
		reducer.Filter(func(b Book) bool { return b.Year >= 2000 },
			reducer.GroupBy(func(b Book) string { return b.Author },
				reducer.Map(func(b Book) string { return b.Title }, reducer.ToSlice[string]()))))

	require.Equal(t, map[string][]string{
		"A": {"S"},
		"B": {"R", "U"},
	}, m)
}

func TestReduce0(t *testing.T) {
	i := iterable.New(100, 20, 3)
	r := reducer.Apply(i, reducer.Reduce0(func(x int, y int) int { return x + y }))
	require.Equal(t, 123, r)
}

func TestCount(t *testing.T) {
	count := reducer.Apply(iterable.New(1, 2, 3), reducer.Count[int]())
	require.Equal(t, 3, count)
}

func TestAverage(t *testing.T) {
	avg := reducer.Apply(iterable.New(5, 10, 6), reducer.Average[int]())
	require.Equal(t, 7.0, avg)
}

func TestAverageEmpty(t *testing.T) {
	avg := reducer.Apply(iterable.New[int](), reducer.Average[int]())
	require.Equal(t, 0.0, avg)
}

func TestMax(t *testing.T) {
	avg := reducer.Apply(iterable.New(5, 10, 6), reducer.Max[int]())
	require.Equal(t, 10, avg)
}

func TestMin(t *testing.T) {
	avg := reducer.Apply(iterable.New(5, 10, 3, 6), reducer.Min[int]())
	require.Equal(t, 3, avg)
}

func TestForallTrue(t *testing.T) {
	require.True(t, reducer.Apply(iterable.New(1, 2, 3, 4), reducer.Forall(func(x int) bool { return x <= 4 })))
}

func TestForallFalse(t *testing.T) {
	require.False(t, reducer.Apply(iterable.New(1, 2, 3, 4), reducer.Forall(func(x int) bool { return x <= 3 })))
}

func TestExistsTrue(t *testing.T) {
	require.True(t, reducer.Apply(iterable.New(1, 2, 3, 4), reducer.Exists(func(x int) bool { return x == 3 })))
}

func TestExistsFalse(t *testing.T) {
	require.False(t, reducer.Apply(iterable.New(1, 2, 3, 4), reducer.Exists(func(x int) bool { return x > 4 })))
}

func TestDoErr(t *testing.T) {
	count := 0
	e := fmt.Errorf("test-error")
	err := reducer.Apply(iterable.New(1, 2, 3, 4), reducer.DoErr(func(x int) error {
		if x > 2 {
			return e
		}
		count++
		return nil
	}))
	require.Equal(t, e, err)
	require.Equal(t, 2, count)
}

func TestDo(t *testing.T) {
	sum := 0
	reducer.Apply(iterable.New(1, 2, 3, 4), reducer.Do(func(x int) {
		sum += x
	}))
	require.Equal(t, 10, sum)
}

func TestGroupByCollect(t *testing.T) {
	m := reducer.Apply(iterable.New(1, 2, 3, 4, 5, 6), reducer.GroupByCollect(func(x int) int { return x % 2 }))
	require.Equal(t, map[int][]int{0: {2, 4, 6}, 1: {1, 3, 5}}, m)
}

func TestToMap(t *testing.T) {
	m := reducer.Apply(iterable.New(1, 2, 3, 4, 5, 6), reducer.ToMap(func(x int) int { return x % 2 }, func(x int) int { return 10 * x }))
	require.Equal(t, map[int]int{0: 20, 1: 10}, m)
}

func TestToMapId(t *testing.T) {
	m := reducer.Apply(iterable.New(1, 2, 3, 4, 5, 6), reducer.ToMapId(func(x int) int { return x % 2 }))
	require.Equal(t, map[int]int{0: 2, 1: 1}, m)
}

func TestToSet(t *testing.T) {
	s := reducer.Apply(iterable.New(1, 2, 3), reducer.ToSet[int]())
	require.Equal(t, map[int]bool{1: true, 2: true, 3: true}, s)
}

func TestApplyMethod(t *testing.T) {
	require.Equal(t, 5, reducer.Max[int]().Apply(iterable.New(1, 5, 3)))
}

func TestApplyIteratorMethod(t *testing.T) {
	require.Equal(t, 5, reducer.Limit(3, reducer.Max[int]()).ApplyIterator(iterable.New(1, 5, 3, 10).Iterator()))
}
