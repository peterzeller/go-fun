package list_test

import (
	"fmt"
	"strconv"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/list/list"
)

func ExampleList_At() {
	l := list.New('a', 'b', 'c')
	fmt.Printf("elem = %c", l.At(1))
	// output: elem = b
}

func ExampleList_Iterator() {
	l := list.New(7, 8, 9)
	sum := 0
	for it := iterable.Start[int](l); it.HasNext(); it.Next() {
		sum += it.Current()
	}
	fmt.Printf("sum = %d", sum)
	// output: sum = 24
}

func ExampleList_Length() {
	l := list.New('a', 'b', 'c')
	fmt.Printf("len = %d", l.Length())
	// output: len = 3
}

func ExampleList_Append() {
	a := list.New(5, 6, 7)
	b := list.New(8, 9)
	c := a.Append(b)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("c = %v\n", c)
	// output: a = [5, 6, 7]
	// b = [8, 9]
	// c = [5, 6, 7, 8, 9]
}

func ExampleList_AppendElems() {
	a := list.New(5, 6, 7)
	b := []int{8, 9}
	c := a.AppendElems(b...).AppendElems(10, 11)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("c = %v\n", c)
	// output: a = [5, 6, 7]
	// b = [8 9]
	// c = [5, 6, 7, 8, 9, 10, 11]
}

func ExampleList_Contains() {
	l := list.New(5, 10, 15)
	fmt.Printf("contains 10: %v\n", l.Contains(10, equality.Default[int]()))
	fmt.Printf("contains 11: %v\n", l.Contains(11, equality.Default[int]()))
	// output: contains 10: true
	// contains 11: false
}

func ExampleList_Equal() {
	a := list.New(5, 6, 7)
	b := list.New(5, 6, 7)
	c := list.New(5, 7, 6)
	fmt.Printf("a equals b: %v\n", a.Equal(b, equality.Default[int]()))
	fmt.Printf("a equals c: %v\n", a.Equal(c, equality.Default[int]()))
	// output: a equals b: true
	// a equals c: false
}

func ExampleList_PrefixOf() {
	a := list.New(5, 6, 7, 8, 9)
	b := list.New(5, 6, 7)
	c := list.New(5, 7, 8)
	fmt.Printf("a prefix of a: %v\n", a.PrefixOf(a, equality.Default[int]()))
	fmt.Printf("b prefix of a: %v\n", b.PrefixOf(a, equality.Default[int]()))
	fmt.Printf("c prefix of a: %v\n", c.PrefixOf(a, equality.Default[int]()))
	// output: a prefix of a: true
	// b prefix of a: true
	// c prefix of a: false
}

func ExampleList_Forall() {
	l := list.New(5, 6, 7)
	allSmallerTen := l.Forall(func(x int) bool { return x < 10 })
	allOdd := l.Forall(func(x int) bool { return x%2 == 1 })
	fmt.Printf("allSmallerTen: %v\n", allSmallerTen)
	fmt.Printf("allOdd: %v\n", allOdd)
	// output: allSmallerTen: true
	// allOdd: false
}

func ExampleList_Exists() {
	l := list.New(10, 122, 42)
	hasGreater100 := l.Exists(func(x int) bool { return x > 100 })
	hasOdd := l.Exists(func(x int) bool { return x%2 == 1 })
	fmt.Printf("hasGreater100: %v\n", hasGreater100)
	fmt.Printf("hasOdd: %v\n", hasOdd)
	// output: hasGreater100: true
	// hasOdd: false
}

func ExampleList_Skip() {
	a := list.New(5, 3, 9, 42, 14)
	b := a.Skip(2)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output: a = [5, 3, 9, 42, 14]
	// b = [9, 42, 14]
}

func ExampleList_Limit() {
	a := list.New(5, 3, 9, 42, 14)
	b := a.Limit(2)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output: a = [5, 3, 9, 42, 14]
	// b = [5, 3]
}

func ExampleList_RemoveAt() {
	a := list.New(5, 3, 9, 42, 14)
	b := a.RemoveAt(2)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output: a = [5, 3, 9, 42, 14]
	// b = [5, 3, 42, 14]
}

func ExampleList_RemoveFirst() {
	a := list.New(5, 3, 3, 9, 42, 14, 3)
	b := a.RemoveFirst(3, equality.Default[int]())
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output: a = [5, 3, 3, 9, 42, 14, 3]
	// b = [5, 3, 9, 42, 14, 3]
}

func ExampleList_RemoveAll() {
	a := list.New(5, 3, 3, 9, 42, 14, 3)
	b := a.RemoveAll(3, equality.Default[int]())
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output: a = [5, 3, 3, 9, 42, 14, 3]
	// b = [5, 9, 42, 14]
}

func ExampleFromIterable() {
	l := list.FromIterable(iterable.RangeI(1, 5))
	fmt.Printf("l = %v\n", l)
	// output: l = [1, 2, 3, 4, 5]
}

func ExampleMap() {
	a := list.New("5", "3", "4")
	b := list.Map(a, func(a string) int {
		b, _ := strconv.Atoi(a)
		return b
	})
	fmt.Printf("b = %v\n", b)
	// output: b = [5, 3, 4]
}

func ExampleMapErr() {
	a := list.New("5", "3", "4")
	b, err := list.MapErr(a, strconv.Atoi)
	fmt.Printf("b = %v, err = %v\n", b, err)
	// output: b = [5, 3, 4], err = <nil>
}

func ExampleMapErr_withError() {
	a := list.New("5", "three", "4")
	b, err := list.MapErr(a, strconv.Atoi)
	fmt.Printf("b = %v, err = %v\n", b, err)
	// output: b = [], err = at index 1: strconv.Atoi: parsing "three": invalid syntax
}

func ExampleFlatMap() {
	a := list.New(1, 2, 3)
	b := list.FlatMap(a, func(a int) iterable.Iterable[int] {
		return list.New(-a, a)
	})
	fmt.Printf("b = %v\n", b)
	// output: b = [-1, 1, -2, 2, -3, 3]
}

func ExampleFilter() {
	a := list.New(7, 9, 10, 8, 11)
	b := a.Filter(func(a int) bool {
		return a >= 9
	})
	fmt.Printf("b = %v\n", b)
	// output: b = [9, 10, 11]
}
