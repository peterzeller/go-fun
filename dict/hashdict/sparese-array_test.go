package hashdict

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/dict"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func genEntries(t *rapid.T, name string) []dict.Entry[int, string] {
	return rapid.SliceOf(rapid.Custom(func(t *rapid.T) dict.Entry[int, string] {
		return dict.Entry[int, string]{
			Key:   rapid.IntRange(0, 31).Draw(t, "key").(int),
			Value: fmt.Sprintf("%d", rapid.IntRange(0, 10).Draw(t, "value").(int)),
		}
	})).Draw(t, name).([]dict.Entry[int, string])
}

func TestSparseArrayGet(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		ar := genEntries(t, "slice")
		values := make([]dict.Entry[int, string], 0, len(ar))
		m := make(map[int]string)

		for _, e := range ar {
			if _, ok := m[e.Key]; !ok {
				values = append(values, e)
				m[e.Key] = e.Value
			}
		}
		sparse := newSparseArray(values...)
		t.Logf("values = %+v", values)
		for i := 0; i < 32; i++ {
			v, ok := sparse.get(i)
			v2, ok2 := m[i]
			require.Equal(t, ok2, ok)
			if ok {
				require.Equal(t, v2, v)
			}
		}
	})
}

func TestSparseArraySet(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		ar := genEntries(t, "slice")
		values := make([]dict.Entry[int, string], 0, len(ar))
		m := make(map[int]string)

		for _, e := range ar {
			if _, ok := m[e.Key]; !ok {
				values = append(values, e)
				m[e.Key] = e.Value
			}
		}

		updates := genEntries(t, "updates")

		t.Logf("values = %+v", values)
		t.Logf("updates = %+v", updates)

		sparse := newSparseArray(values...)
		for _, e := range updates {
			sparse = sparse.set(e.Key, e.Value)
			m[e.Key] = e.Value
		}

		for i := 0; i < 32; i++ {
			v, ok := sparse.get(i)
			v2, ok2 := m[i]
			require.Equal(t, ok2, ok)
			if ok {
				require.Equal(t, v2, v)
			}
		}
	})
}

func TestSparseArrayRemove(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		ar := genEntries(t, "slice")
		values := make([]dict.Entry[int, string], 0, len(ar))
		m := make(map[int]string)

		for _, e := range ar {
			if _, ok := m[e.Key]; !ok {
				values = append(values, e)
				m[e.Key] = e.Value
			}
		}

		updates := rapid.SliceOf(rapid.IntRange(0, 31)).Draw(t, "removes").([]int)

		t.Logf("values = %+v", values)
		t.Logf("updates = %+v", updates)

		sparse := newSparseArray(values...)
		for _, i := range updates {
			sparse = sparse.remove(i)
			delete(m, i)
		}

		for i := 0; i < 32; i++ {
			v, ok := sparse.get(i)
			v2, ok2 := m[i]
			require.Equal(t, ok2, ok)
			if ok {
				require.Equal(t, v2, v)
			}
		}
	})
}
