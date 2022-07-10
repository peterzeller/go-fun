package linked_test

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/list/linked"
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

func ExampleList_Reversed() {
	list := linked.New(1, 2, 3, 4, 5, 6, 7, 8)
	reversed := list.Reversed()
	fmt.Printf("list = %v\n", list)
	fmt.Printf("reversed = %v\n", reversed)
	// output: list = [1, 2, 3, 4, 5, 6, 7, 8]
	// reversed = [8, 7, 6, 5, 4, 3, 2, 1]
}

func ExampleList_Iterator() {
	list := linked.New(1, 2, 3)
	it := list.Iterator()
	for {
		v, ok := it.Next()
		if !ok {
			break
		}
		fmt.Printf("%d, ", v)
	}
	// output: 1, 2, 3,
}

func ExampleList_Head() {
	list := linked.New(7, 8, 9)
	fmt.Printf("head = %d", list.Head())
	// output: head = 7
}

func ExampleList_Tail() {
	list := linked.New(7, 8, 9)
	tail := list.Tail()
	fmt.Printf("list = %v\n", list)
	fmt.Printf("tail = %v\n", tail)
	// output: list = [7, 8, 9]
	// tail = [8, 9]
}

func ExampleList_Contains() {
	list := linked.New("banana", "apple", "orange")
	fmt.Printf("contains apple: %v\n", list.Contains("apple", equality.Default[string]()))
	fmt.Printf("contains tomato: %v\n", list.Contains("tomato", equality.Default[string]()))
	// output: contains apple: true
	// contains tomato: false
}

func ExampleList_Equal() {
	a := linked.New("banana", "apple", "orange")
	b := linked.New("banana", "apple", "orange")
	c := linked.New("banana", "orange", "apple")
	fmt.Printf("a equals b: %v\n", a.Equal(b, equality.Default[string]()))
	fmt.Printf("a equals c: %v\n", a.Equal(c, equality.Default[string]()))
	// output: a equals b: true
	// a equals c: false
}

func ExampleList_PrefixOf() {
	a := linked.New(7, 8, 9, 10, 11)
	b := linked.New(7, 8, 9)
	c := linked.New(7, 9, 10)
	fmt.Printf("a prefix of a: %v\n", a.PrefixOf(a, equality.Default[int]()))
	fmt.Printf("b prefix of a: %v\n", b.PrefixOf(a, equality.Default[int]()))
	fmt.Printf("c prefix of a: %v\n", c.PrefixOf(a, equality.Default[int]()))
	// output: a prefix of a: true
	// b prefix of a: true
	// c prefix of a: false
}

func ExampleList_Forall() {
	l := linked.New(5, 6, 7)
	allSmallerTen := l.Forall(func(x int) bool { return x < 10 })
	allOdd := l.Forall(func(x int) bool { return x%2 == 1 })
	fmt.Printf("allSmallerTen: %v\n", allSmallerTen)
	fmt.Printf("allOdd: %v\n", allOdd)
	// output: allSmallerTen: true
	// allOdd: false
}

func ExampleList_Exists() {
	l := linked.New(10, 122, 42)
	hasGreater100 := l.Exists(func(x int) bool { return x > 100 })
	hasOdd := l.Exists(func(x int) bool { return x%2 == 1 })
	fmt.Printf("hasGreater100: %v\n", hasGreater100)
	fmt.Printf("hasOdd: %v\n", hasOdd)
	// output: hasGreater100: true
	// hasOdd: false
}

func ExampleList_Skip() {
	a := linked.New(5, 3, 9, 42, 14)
	b := a.Skip(2)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output: a = [5, 3, 9, 42, 14]
	// b = [9, 42, 14]
}

func ExampleList_Limit() {
	a := linked.New(5, 3, 9, 42, 14)
	b := a.Limit(2)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output: a = [5, 3, 9, 42, 14]
	// b = [5, 3]
}
