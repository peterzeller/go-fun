package hashdict_test

import (
	"testing"

	"github.com/peterzeller/go-fun/dict"
	"github.com/peterzeller/go-fun/dict/hashdict"
	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/list"
	"github.com/peterzeller/go-fun/reducer"
	"github.com/stretchr/testify/assert"
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

func TestGetOr(t *testing.T) {
	d := hashdict.New(hash.String(), dict.E("a", 1), dict.E("b", 2))
	require.Equal(t, 1, d.GetOrZero("a"))
	require.Equal(t, 2, d.GetOr("b", 42))
	require.Equal(t, 0, d.GetOrZero("xyz"))
	require.Equal(t, 42, d.GetOr("xyz", 42))
}

func TestKeys(t *testing.T) {
	d := hashdict.New(hash.String(), dict.E("a", 1), dict.E("b", 2), dict.E("c", 3))

	strings := reducer.Apply(d.Keys(), reducer.ToSet[string]())

	expected := map[string]bool{
		"a": true,
		"b": true,
		"c": true,
	}

	require.Equal(t, expected, strings)
}

func TestValues(t *testing.T) {
	d := hashdict.New(hash.String(), dict.E("a", 1), dict.E("b", 2), dict.E("c", 3))

	strings := reducer.Apply(d.Values(), reducer.ToSet[int]())

	expected := map[int]bool{
		1: true,
		2: true,
		3: true,
	}

	require.Equal(t, expected, strings)
}

func TestMap(t *testing.T) {
	d := hashdict.New(hash.String(), dict.E("a", 1), dict.E("b", 2), dict.E("c", 3))

	d = d.Map(func(key string, value int) int { return value * 10 })

	require.Equal(t, 10, d.GetOrZero("a"))
	require.Equal(t, 20, d.GetOrZero("b"))
	require.Equal(t, 30, d.GetOrZero("c"))
}

func TestFilter(t *testing.T) {
	d := hashdict.New(hash.String(), dict.E("a", 1), dict.E("b", 2), dict.E("c", 3))

	d = d.Filter(func(key string, value int) bool { return value%2 == 1 })

	require.Equal(t, 2, d.Size())
	require.Equal(t, 1, d.GetOrZero("a"))
	require.Equal(t, 0, d.GetOrZero("b"))
	require.Equal(t, 3, d.GetOrZero("c"))
}

func TestFromMap(t *testing.T) {
	a := hashdict.New(hash.String(), dict.E("a", 1), dict.E("b", 2), dict.E("c", 3))
	b := hashdict.FromMap(hash.String(), map[string]int{"a": 1, "b": 2, "c": 3})
	require.True(t, a.Equal(b, equality.Default[int]()))
}

func TestMergeAll(t *testing.T) {
	a := hashdict.New(hash.String(), dict.E("a", 1), dict.E("b", 2), dict.E("c", 3))

	b := a.MergeAll(list.New(dict.E("b", 20), dict.E("d", 40)), hashdict.MergeOpts[string, int, int, int]{
		Left:  func(k string, a int) (int, bool) { return a, true },
		Right: func(k string, b int) (int, bool) { return b, true },
		Both:  func(k string, a int, b int) (int, bool) { return b, true },
	})

	expected := hashdict.FromMap(hash.String(), map[string]int{"a": 1, "b": 20, "c": 3, "d": 40})

	require.True(t, b.Equal(expected, equality.Default[int]()))
}

func TestEqual(t *testing.T) {
	a := hashdict.FromMap(hash.String(), map[string]int{"a": 1, "b": 2, "c": 3})
	b := hashdict.FromMap(hash.String(), map[string]int{"c": 3})
	c := hashdict.FromMap(hash.String(), map[string]int{"a": 1, "b": 20, "c": 3})

	assert.True(t, a.Equal(a, equality.Default[int]()))
	assert.False(t, a.Equal(b, equality.Default[int]()))
	assert.False(t, a.Equal(c, equality.Default[int]()))
	assert.False(t, b.Equal(a, equality.Default[int]()))
	assert.True(t, b.Equal(b, equality.Default[int]()))
	assert.False(t, b.Equal(c, equality.Default[int]()))
	assert.False(t, c.Equal(a, equality.Default[int]()))
	assert.False(t, c.Equal(b, equality.Default[int]()))
	assert.True(t, c.Equal(c, equality.Default[int]()))
}
