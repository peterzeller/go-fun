package reducer_test

import (
	"testing"

	"github.com/peterzeller/go-fun/v2/iterable"
	"github.com/peterzeller/go-fun/v2/reducer"
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
