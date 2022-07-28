# Go-fun

![pipeline status](https://github.com/peterzeller/go-fun/actions/workflows/go.yml/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/peterzeller/go-fun/badge.svg)](https://coveralls.io/github/peterzeller/go-fun) [![Go Reference](https://pkg.go.dev/badge/github.com/peterzeller/go-fun.svg)](https://pkg.go.dev/github.com/peterzeller/go-fun)

Utilities and immutable collections for functional programming in Golang.

For documentation, check the [Go Package Documentation](https://pkg.go.dev/github.com/peterzeller/go-fun).

## Features

- Immutable data structures
    - List (package [list](./list))
      - Slice based (package [list](./list/list))
      - Singly linked list (package [linked](./list/linked))
    - Dict (package [dict](./dict))
        - HashDict (package [dict/hashdict](./dict/hashdict))
        - ArrayDict (package [dict/arraydict](./dict/arraydict))
    - Set (package [set](./set))
        - HashSet (package [dict/hashset](./dict/hashset))
    - Optional (package [opt](./opt))
- Iterable abstraction (package [iterable](./iterable))
- Reducers for transforming data (map, filter, group by, etc) (package [reducer](./reducer))
- Equality type class (package [equality](./equality))
- Hash type class (package [hash](./hash))
- Generic Zero Value (package [zero](./zero))
- Generic Slice functions (package [slice](./slice))
- Mutable data structures (package [mutable](./mutable))
    - Stack 

## Why immutable collections?

- Immutable data is thread-safe by default
- Immutable data structures are more efficient than copying data to prevent unwanted modification
- Use immutable data structures at API boundaries to make it clear that data cannot be modified.
- Updated versions of a data structure share underlying state, which makes them more memory efficient when keeping multiple versions (in recursive algorithms, in search algorithms, for undo/history functionality)


## Examples

### Immutable Dict

	d0 := hashdict.New[string, int](hash.String())

	d1 := d0.Set("a", 1)
	d2 := d1.Set("b", 42)
	d3 := d2.Set("a", 7)

	require.Equal(t, 1, d1.GetOrZero("a"))
	require.Equal(t, 1, d2.GetOrZero("a"))
	require.Equal(t, 7, d3.GetOrZero("a"))

	require.Equal(t, 0, d1.GetOrZero("b"))
	require.Equal(t, 42, d2.GetOrZero("b"))
	require.Equal(t, 42, d3.GetOrZero("b"))


### Immutable Set


### Reducers

Take a slice of books and return a map where the keys are authors and the values contain all the titles of the books the author has written in or after the year 2000.


```
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
```
