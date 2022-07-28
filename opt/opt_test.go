package opt_test

import (
	"encoding/json"
	"fmt"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/list/list"
	"github.com/peterzeller/go-fun/opt"
)

func ExampleFirst() {
	list1 := list.New(7, 8, 9)
	list2 := list.New[int]()
	fmt.Println("first of list1 =", opt.First[int](list1))
	fmt.Println("first of list2 =", opt.First[int](list2))
	// output:
	// first of list1 = Some(7)
	// first of list2 = None()
}

func ExampleOptional_Get() {
	x := opt.Some(42)
	if v, ok := x.Get(); ok {
		fmt.Println("value in x is", v)
	} else {
		fmt.Println("no value in x")
	}
	// output:
	// value in x is 42
}

func ExampleOptional_Iterator() {
	x := opt.Some(5)
	y := opt.First(iterable.Map[int, string](x, func(i int) string {
		return fmt.Sprintf("%d", i*i)
	}))
	fmt.Println("x =", x)
	fmt.Println("y =", y)
	// output:
	// x = Some(5)
	// y = Some(25)
}

type Pair struct {
	X opt.Optional[int]
	Y opt.Optional[int]
}

func ExampleOptional_MarshalJSON() {
	example := Pair{
		X: opt.Some(42),
		Y: opt.None[int](),
	}
	data, err := json.Marshal(example)
	if err != nil {
		panic(err)
	}
	fmt.Println("data =", string(data))
	// output:
	// data = {"X":42,"Y":null}
}

func ExampleOptional_UnmarshalJSON() {
	data := []byte(`{"X":42,"Y":null}`)
	var example Pair
	err := json.Unmarshal(data, &example)
	if err != nil {
		panic(err)
	}
	fmt.Printf("example = %+v\n", example)
	// output:
	// example = {X:Some(42) Y:None()}
}

func ExampleOptional_OrElse() {
	a := opt.None[int]()
	b := opt.Some(42)
	fmt.Println("a ->", a.OrElse(100))
	fmt.Println("b ->", b.OrElse(100))
	// output:
	// a -> 100
	// b -> 42
}

func ExampleOptional_OrElseGet() {
	a := opt.None[int]()
	b := opt.Some(42)
	fmt.Println("a ->", a.OrElseGet(func() int { return 100 }))
	fmt.Println("b ->", b.OrElseGet(func() int { return 100 }))
	// output:
	// a -> 100
	// b -> 42
}

func ExampleOptional_OrElsePanic() {
	x := opt.Some(42)
	fmt.Println("x ->", x.OrElsePanic())
	// output:
	// x -> 42
}

func ExampleOptional_Present() {
	a := opt.None[int]()
	b := opt.Some(42)
	fmt.Println("a ->", a.Present())
	fmt.Println("b ->", b.Present())
	// output:
	// a -> false
	// b -> true
}

func ExampleOptional_String() {
	a := opt.None[int]()
	b := opt.Some(42)
	fmt.Println("a ->", a.String())
	fmt.Println("b ->", b.String())
	// output:
	// a -> None()
	// b -> Some(42)
}

func ExampleOptional_Value() {
	a := opt.None[int]()
	b := opt.Some(42)
	fmt.Println("a ->", a.Value())
	fmt.Println("b ->", b.Value())
	// output:
	// a -> 0
	// b -> 42
}

func ExampleOptional_GetPointer() {
	a := opt.None[int]()
	b := opt.Some(42)
	fmt.Printf("a -> %#v\n", a.GetPointer())
	fmt.Printf("b -> %#v\n", *b.GetPointer())
	// output:
	// a -> (*int)(nil)
	// b -> 42
}
