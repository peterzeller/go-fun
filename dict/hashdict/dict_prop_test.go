package hashdict

import (
	"fmt"
	"log"
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

func TestDicSetMerge(t *testing.T) {
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
			case 0: // set
				key := genKey(t)
				v := rapid.IntRange(0, 10).Draw(t, "value").(int)
				t.Logf("dict[%d] = dict[%d].Set('%s', %d)", len(dicts), d, key, v)
				d1 := dict.Set(key, v)
				d2 := arrayDict.Set(key, v, keyHash)
				dicts = append(dicts, d1)
				arrayDicts = append(arrayDicts, d2)
			case 1: // merge left
				left := rapid.IntRange(0, len(dicts)-1).Draw(t, "left").(int)
				right := rapid.IntRange(0, len(dicts)-1).Draw(t, "right").(int)
				t.Logf("dict[%d] = dict[%d].MergeLeft(dict[%d])", len(dicts), left, right)
				d1 := dicts[left].MergeLeft(dicts[right])
				d2 := arrayDicts[left].MergeLeft(arrayDicts[right], keyHash)
				dicts = append(dicts, d1)
				arrayDicts = append(arrayDicts, d2)
			case 2: // merge right
				left := rapid.IntRange(0, len(dicts)-1).Draw(t, "left").(int)
				right := rapid.IntRange(0, len(dicts)-1).Draw(t, "right").(int)
				t.Logf("dict[%d] = dict[%d].MergeRight(dict[%d])", len(dicts), left, right)
				d1 := dicts[left].MergeLeft(dicts[right])
				d2 := arrayDicts[left].MergeLeft(arrayDicts[right], keyHash)
				dicts = append(dicts, d1)
				arrayDicts = append(arrayDicts, d2)
			}
			require.NoError(t, dicts[len(dicts)-1].checkInvariant())
			for i := range dicts {
				assertDictsEqual(t, arrayDicts[i], dicts[i])
			}
		}
	})
}

func assertDictsEqual(t require.TestingT, a arraydict.ArrayDict[key, int], b Dict[key, int]) {
	var err error
	defer func() {
		if err != nil {
			require.NoError(t, err, "maps differ:\na = %s\nb = %s", a.String(), b.String())
		}
	}()
	for it := a.Iterator(); ; {
		ae, ok := it.Next()
		if !ok {
			break
		}
		bv, ok := b.Get(ae.Key)
		if !ok {
			err = fmt.Errorf("Key %+v exists in b but not in a", ae.Key)
			return
		}
		if ae.Value != bv {
			err = fmt.Errorf("values differ for key %+v : %+v / %+v", ae.Key, ae.Value, bv)
			return
		}
	}
	for it := b.Iterator(); ; {
		be, ok := it.Next()
		if !ok {
			break
		}
		av, ok := a.Get(be.Key, keyHash)
		if !ok {
			err = fmt.Errorf("Key %+v exists in a but not in b", be.Key)
			return
		}
		if be.Value != av {
			err = fmt.Errorf("values differ for key %+v : %+v / %+v", be.Key, be.Value, av)
			return
		}
	}
}

func TestMerge(t *testing.T) {
	d0 := New[key, int](keyHash)
	d1 := d0.Set("", 0)
	d2 := d0.Set("b", 0)
	d3 := d1.Set("e", 0)
	log.Printf("d2 = %+v", d2)
	log.Printf("d3 = %+v", d3)
	d4 := d2.MergeRight(d3)
	log.Printf("d4 = %+v", d4)
	bv, ok := d4.Get("b")
	require.True(t, ok)
	require.Equal(t, 0, bv)
}
