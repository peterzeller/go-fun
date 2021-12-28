package hashdict

import (
	"strings"
	"testing"

	"github.com/peterzeller/go-fun/v2/dict/arraydict"
	"github.com/peterzeller/go-fun/v2/hash"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

type key string

var keyHash hash.EqHash[key] = hash.Fun[key]{
	Eq: func(a, b key) bool {
		return a == b
	},
	// hash function with lots of collisions to make it interesting
	H: func(a key) (h int64) {
		for i := 0; i < len(a); i++ {
			h = 2 * h
			if a[i] == 'a' {
				h++
			}
		}
		return
	},
}

func genKey(t *rapid.T) key {
	return rapid.SliceOf(rapid.IntRange(0, 4)).
		Map(func(s []int) key {
			var res strings.Builder
			for _, c := range s {
				res.WriteRune(rune('a' + c))
			}
			return key(res.String())
		}).
		Draw(t, "key").(key)

}

func TestDictGetSet(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {

		dicts := []Dict[key, int]{
			New[key, int](keyHash),
		}
		arrayDicts := []arraydict.ArrayDict[key, int]{
			arraydict.New[key, int](),
		}

		n := rapid.IntRange(1, 100).Draw(t, "n").(int)
		for i := 0; i < n; i++ {
			d := rapid.IntRange(0, len(dicts)-1).Draw(t, "d").(int)
			dict := dicts[d]
			arrayDict := arrayDicts[d]
			cmd := rapid.IntRange(0, 1).Draw(t, "cmd").(int)
			switch cmd {
			case 0: // get
				key := genKey(t)
				t.Logf("dict[%d].Get('%s')", d, key)
				v1, ok1 := dict.Get(key)
				v2, ok2 := arrayDict.Get(key, keyHash)
				require.Equal(t, ok2, ok1)
				require.Equal(t, v2, v1)
			case 1: // set
				key := genKey(t)
				v := rapid.IntRange(0, 10).Draw(t, "value").(int)
				t.Logf("dict[%d] = dict[%d].Set('%s', %d)", len(dicts), d, key, v)
				d1 := dict.Set(key, v)
				d2 := arrayDict.Set(key, v, keyHash)
				dicts = append(dicts, d1)
				arrayDicts = append(arrayDicts, d2)
			}
			for i := range dicts {
				assertDictsEqual(t, arrayDicts[i], dicts[i])
			}
			require.NoError(t, dicts[len(dicts)-1].checkInvariant())
		}
	})
}

func TestDictGetSetRemove(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {

		dicts := []Dict[key, int]{
			New[key, int](keyHash),
		}
		arrayDicts := []arraydict.ArrayDict[key, int]{
			arraydict.New[key, int](),
		}

		n := rapid.IntRange(1, 100).Draw(t, "n").(int)
		for i := 0; i < n; i++ {
			d := rapid.IntRange(0, len(dicts)-1).Draw(t, "d").(int)
			dict := dicts[d]
			arrayDict := arrayDicts[d]
			cmd := rapid.IntRange(0, 2).Draw(t, "cmd").(int)
			switch cmd {
			case 0: // get
				key := genKey(t)
				t.Logf("dict[%d].Get('%s')", d, key)
				v1, ok1 := dict.Get(key)
				v2, ok2 := arrayDict.Get(key, keyHash)
				require.Equal(t, ok2, ok1)
				require.Equal(t, v2, v1)
			case 1: // set
				key := genKey(t)
				v := rapid.IntRange(0, 10).Draw(t, "value").(int)
				t.Logf("dict[%d] = dict[%d].Set('%s', %d)", len(dicts), d, key, v)
				d1 := dict.Set(key, v)
				d2 := arrayDict.Set(key, v, keyHash)
				dicts = append(dicts, d1)
				arrayDicts = append(arrayDicts, d2)
			case 2: // remove
				key := genKey(t)
				t.Logf("dict[%d] = dict[%d].Remove('%s')", len(dicts), d, key)
				d1 := dict.Remove(key)
				d2, _ := arrayDict.Remove(key, keyHash)
				dicts = append(dicts, d1)
				arrayDicts = append(arrayDicts, d2)
			}
			for i := range dicts {
				assertDictsEqual(t, arrayDicts[i], dicts[i])
			}
			require.NoError(t, dicts[len(dicts)-1].checkInvariant())
		}
	})
}

func assertDictsEqual(t require.TestingT, a arraydict.ArrayDict[key, int], b Dict[key, int]) {
	for it := a.Iterator(); ; {
		ae, ok := it.Next()
		if !ok {
			break
		}
		bv, ok := b.Get(ae.Key)
		require.True(t, ok, "Key %+v exists in b but not in a", ae.Key)
		require.Equal(t, ae.Value, bv)
	}
	for it := b.Iterator(); ; {
		be, ok := it.Next()
		if !ok {
			break
		}
		av, ok := a.Get(be.Key, keyHash)
		require.True(t, ok, "Key %+v exists in a but not in b", be.Key)
		require.Equal(t, be.Value, av, "Key %+v has different values for b and a", be.Key)
	}
}
