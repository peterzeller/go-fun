package iterable_test

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
)

func ExampleLength() {
	it := iterable.New(4, 5, 6, 7)
	l := iterable.Length(it)
	fmt.Printf("l = %d\n", l)
	// output: l = 4
}
