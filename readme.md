# Go-fun

![pipeline status](https://github.com/peterzeller/go-fun/actions/workflows/go.yml/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/peterzeller/go-fun/badge.svg)](https://coveralls.io/github/peterzeller/go-fun)

Utilities and immutable collections for functional programming in Golang.

This is an experimental library to play with the new Generics Feature in Go 1.18.


## Features

- Reducers for transforming data (map, filter, group by, etc)
- Immutable data structures
    - List
- Mutable data structures
    - Stack
- Iterable abstraction
- Equality type class

## Examples

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
