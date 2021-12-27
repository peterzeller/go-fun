package hashdict_test

import (
	"testing"

	"github.com/peterzeller/go-fun/v2/dict/hashdict"
	"github.com/peterzeller/go-fun/v2/hash"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
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
}
