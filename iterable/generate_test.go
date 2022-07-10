package iterable_test

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	parts := iterable.TakeWhile(
		func(x int) bool {
			return x > 0
		}, iterable.Generate[int](
			10,
			func(x int) int {
				return x / 2
			}))

	require.Equal(t, []int{10, 5, 2, 1}, iterable.ToSlice(parts))
}

func ExampleGenerate() {
	it := iterable.Generate[int](
		10,
		func(x int) int {
			return x / 2
		})

	fmt.Printf("it = %s", iterable.String(iterable.Take(5, it)))
	// output: it = [10, 5, 2, 1, 0]
}

func ExampleGenerateState() {
	it := iterable.GenerateState[int](
		0,
		func(state int) (newState int, res int, ok bool) {
			ok = state <= 10
			res = 10 * state
			newState = state + 2
			return
		})

	fmt.Printf("it = %s", iterable.String(it))
	// output: it = [0, 20, 40, 60, 80, 100]
}
