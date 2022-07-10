package zero_test

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/zero"
	"github.com/stretchr/testify/require"
)

func TestZero(t *testing.T) {
	require.Equal(t, "", zero.Value[string]())
	require.Equal(t, 0, zero.Value[int]())
}

func ExampleValue() {
	fmt.Printf("zero value is %v\n", zero.Value[int]())
	// output: zero value is 0
}
