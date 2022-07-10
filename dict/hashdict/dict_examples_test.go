package hashdict_test

import (
	"fmt"

	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/dict/hashdict"
	"github.com/peterzeller/go-fun/hash"
)

func ExampleDict_Get() {
	d := hashdict.New(hash.String(),
		dict.E("a", 1),
		dict.E("b", 2),
		dict.E("c", 3),
	)
	v, ok := d.Get("b")
	fmt.Printf("v = %v, ok = %v\n", v, ok)
	// output: v = 2, ok = true
}

func ExampleDict_Get_nonExistent() {
	d := hashdict.New(hash.String(),
		dict.E("a", 1),
		dict.E("b", 2),
		dict.E("c", 3),
	)
	v, ok := d.Get("x")
	fmt.Printf("v = %v, ok = %v\n", v, ok)
	// output: v = 0, ok = false
}

func ExampleDict_ContainsKey() {
	d := hashdict.New(hash.String(),
		dict.E("a", 1),
		dict.E("b", 2),
		dict.E("c", 3),
	)
	fmt.Printf("contains a: %v\n", d.ContainsKey("a"))
	fmt.Printf("contains x: %v\n", d.ContainsKey("x"))
	// output:
	// contains a: true
	// contains x: false
}

func ExampleDict_Set() {
	a := hashdict.New(hash.String(),
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	b := a.Set("x", 42)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output:
	// a = [x -> 1, z -> 3, y -> 2]
	// b = [x -> 42, z -> 3, y -> 2]
}

func ExampleDict_Remove() {
	a := hashdict.New(hash.String(),
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	b := a.Remove("x")
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	// output:
	// a = [x -> 1, z -> 3, y -> 2]
	// b = [z -> 3, y -> 2]
}

func ExampleDict_Size() {
	d := hashdict.New(hash.String(),
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	fmt.Printf("size = %v\n", d.Size())
	// output:
	// size = 3
}

func ExampleDict_Iterator() {
	d := hashdict.New(hash.String(),
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
	// key: z, value: 3
	// key: y, value: 2
}

func ExampleFilterMap() {
	a := hashdict.New(hash.String(),
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	b := hashdict.FilterMap(a, func(key string, val int) ([]int, bool) {
		if key == "y" {
			return nil, false
		}
		return []int{val, 10 * val}, true
	})
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)

	// output:
	// a = [x -> 1, z -> 3, y -> 2]
	// b = [x -> [1 10], z -> [3 30]]
}

func ExampleDict_MergeLeft() {
	a := hashdict.New(hash.String(),
		dict.E("x", 1),
		dict.E("y", 2),
	)
	b := hashdict.New(hash.String(),
		dict.E("y", 20),
		dict.E("z", 30),
	)
	c := a.MergeLeft(b)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("c = %v\n", c)
	// output:
	// a = [x -> 1, y -> 2]
	// b = [z -> 30, y -> 20]
	// c = [x -> 1, z -> 30, y -> 2]
}

func ExampleDict_MergeRight() {
	a := hashdict.New(hash.String(),
		dict.E("x", 1),
		dict.E("y", 2),
	)
	b := hashdict.New(hash.String(),
		dict.E("y", 20),
		dict.E("z", 30),
	)
	c := a.MergeRight(b)
	fmt.Printf("a = %v\n", a)
	fmt.Printf("b = %v\n", b)
	fmt.Printf("c = %v\n", c)
	// output:
	// a = [x -> 1, y -> 2]
	// b = [z -> 30, y -> 20]
	// c = [x -> 1, z -> 30, y -> 20]
}

func ExampleDict_String() {
	d := hashdict.New(hash.String(),
		dict.E("x", 1),
		dict.E("y", 2),
		dict.E("z", 3),
	)
	fmt.Printf("d = %s", d.String())
	// output: d = [x -> 1, z -> 3, y -> 2]
}
