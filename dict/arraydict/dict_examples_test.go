package arraydict_test

import (
	"fmt"

	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/dict/arraydict"
	"github.com/peterzeller/go-fun/equality"
)

func ExampleArrayDict_Get() {
	d := arraydict.New(
		dict.E("a", 1),
		dict.E("b", 2),
		dict.E("c", 3),
	)
	v, ok := d.Get("b", equality.Default[string]())
	fmt.Printf("v = %v, ok = %v\n", v, ok)
	// output: v = 2, ok = true
}

func ExampleArrayDict_Get_nonExistent() {
	d := arraydict.New(
		dict.E("a", 1),
		dict.E("b", 2),
		dict.E("c", 3),
	)
	v, ok := d.Get("x", equality.Default[string]())
	fmt.Printf("v = %v, ok = %v\n", v, ok)
	// output: v = 0, ok = false
}

func ExampleArrayDict_ContainsKey() {
	d := arraydict.New(
		dict.E("a", 1),
		dict.E("b", 2),
		dict.E("c", 3),
	)
	fmt.Printf("contains a: %v\n", d.ContainsKey("a", equality.Default[string]()))
	fmt.Printf("contains x: %v\n", d.ContainsKey("x", equality.Default[string]()))
	// output:
	// contains a: true
	// contains x: false
}

func ExampleArrayDict_Set() {
	a := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	b := a.Set("x", 42, equality.Default[string]())
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output:
	// a = [x -> 1, y -> 2, z -> 3]
	// b = [x -> 42, y -> 2, z -> 3]
}

func ExampleArrayDict_Remove() {
	a := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	b, changed := a.Remove("x", equality.Default[string]())
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("changed = %v\n", changed)
	// output:
	// a = [x -> 1, y -> 2, z -> 3]
	// b = [y -> 2, z -> 3]
	// changed = true
}

func ExampleArrayDict_Size() {
	d := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	fmt.Printf("size = %v\n", d.Size())
	// output:
	// size = 3
}

func ExampleArrayDict_First() {
	d := arraydict.New(
		dict.E("z", 1),
		dict.E("x", 2),
		dict.E("y", 3),
	)
	fmt.Printf("first = %v\n", d.First())
	// output:
	// first = z -> 1
}

func ExampleArrayDict_Iterator() {
	d := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	it := d.Iterator()
	for {
		e, ok := it.Next()
		if !ok {
			break
		}
		fmt.Printf("key: %v, value: %v\n", e.Key, e.Value)
	}
	// output:
	// key: x, value: 1
	// key: y, value: 2
	// key: z, value: 3
}

func ExampleFilterMap() {
	a := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	b := arraydict.FilterMap(a, func(key string, val int) ([]int, bool) {
		if key == "y" {
			return nil, false
		}
		return []int{val, 10 * val}, true
	})
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)

	// output:
	// a = [x -> 1, y -> 2, z -> 3]
	// b = [x -> [1 10], z -> [3 30]]
}

func ExampleArrayDict_MergeLeft() {
	a := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
	)
	b := arraydict.New(
		dict.E("y", 20),
		dict.E("z", 30),
	)
	c := a.MergeLeft(b, equality.Default[string]())
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("c = %v\n", c)
	// output:
	// a = [x -> 1, y -> 2]
	// b = [y -> 20, z -> 30]
	// c = [x -> 1, y -> 2, z -> 30]
}

func ExampleArrayDict_MergeRight() {
	a := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
	)
	b := arraydict.New(
		dict.E("y", 20),
		dict.E("z", 30),
	)
	c := a.MergeRight(b, equality.Default[string]())
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("c = %v\n", c)
	// output:
	// a = [x -> 1, y -> 2]
	// b = [y -> 20, z -> 30]
	// c = [x -> 1, y -> 20, z -> 30]
}

func ExampleArrayDict_String() {
	d := arraydict.New(
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	fmt.Printf("d = %s", d.String())
	// output: d = [x -> 1, y -> 2, z -> 3]
}
