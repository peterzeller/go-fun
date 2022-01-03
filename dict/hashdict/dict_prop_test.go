package hashdict

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/peterzeller/go-fun/v2/dict"
	"github.com/peterzeller/go-fun/v2/dict/arraydict"
	"github.com/peterzeller/go-fun/v2/hash"
	"github.com/peterzeller/go-fun/v2/iterable"
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

func genEntry() *rapid.Generator { // dict.Entry[key, int]
	return rapid.Custom(func(t *rapid.T) dict.Entry[key, int] {
		return dict.Entry[key, int]{
			Key:   genKey(t),
			Value: rapid.IntRange(0, 10).Draw(t, "value").(int),
		}
	})
}

func genDict() *rapid.Generator { // Dict[key, int]
	return rapid.SliceOf(genEntry()).Map(func(ar []dict.Entry[key, int]) Dict[key, int] {
		return New(keyHash, ar...)
	})
}

func TestMergeLeft(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		a := genDict().Draw(t, "a").(Dict[key, int])
		b := genDict().Draw(t, "b").(Dict[key, int])
		t.Logf("a = %+v", a)
		require.NoError(t, a.checkInvariant())
		t.Logf("b = %+v", b)
		require.NoError(t, b.checkInvariant())
		c := a.MergeLeft(b)
		t.Logf("c = %+v", c)
		require.NoError(t, c.checkInvariant())

		for it := iterable.Start[dict.Entry[key, int]](a); it.HasNext(); it.Next() {
			require.Equal(t, it.Current().Value, c.GetOrZero(it.Current().Key), "left value for key %+v", it.Current().Key)
		}

		for it := iterable.Start[dict.Entry[key, int]](b); it.HasNext(); it.Next() {
			if !a.ContainsKey(it.Current().Key) {
				require.Equal(t, it.Current().Value, c.GetOrZero(it.Current().Key), "right value for key %+v", it.Current().Key)
			}
		}
	})
}

func TestMergeRight(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		a := genDict().Draw(t, "a").(Dict[key, int])
		b := genDict().Draw(t, "b").(Dict[key, int])
		t.Logf("a = %+v", a)
		require.NoError(t, a.checkInvariant())
		t.Logf("b = %+v", b)
		require.NoError(t, b.checkInvariant())
		c := a.MergeRight(b)
		t.Logf("c = %+v", c)
		require.NoError(t, c.checkInvariant())

		for it := iterable.Start[dict.Entry[key, int]](b); it.HasNext(); it.Next() {
			require.Equal(t, it.Current().Value, c.GetOrZero(it.Current().Key), "left value for key %+v", it.Current().Key)
		}

		for it := iterable.Start[dict.Entry[key, int]](a); it.HasNext(); it.Next() {
			if !a.ContainsKey(it.Current().Key) {
				require.Equal(t, it.Current().Value, c.GetOrZero(it.Current().Key), "right value for key %+v", it.Current().Key)
			}
		}
	})
}

func TestMergeLeftIterable(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		a := genDict().Draw(t, "a").(Dict[key, int])
		b := rapid.SliceOf(genEntry()).Draw(t, "b").([]dict.Entry[key, int])
		t.Logf("a = %+v", a)
		require.NoError(t, a.checkInvariant())
		t.Logf("b = %+v", b)
		c := a.MergeLeft(iterable.FromSlice(b))
		t.Logf("c = %+v", c)
		require.NoError(t, c.checkInvariant())

		for it := iterable.Start[dict.Entry[key, int]](a); it.HasNext(); it.Next() {
			require.Equal(t, it.Current().Value, c.GetOrZero(it.Current().Key), "left value for key %+v", it.Current().Key)
		}

		for _, e := range b {
			if !a.ContainsKey(e.Key) {
				require.Equal(t, e.Value, c.GetOrZero(e.Key), "right value for key %+v", e.Key)
			}
		}
	})
}

func TestMergeRightIterable(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		a := genDict().Draw(t, "a").(Dict[key, int])
		b := rapid.SliceOf(genEntry()).Draw(t, "b").([]dict.Entry[key, int])
		t.Logf("a = %+v", a)
		require.NoError(t, a.checkInvariant())
		t.Logf("b = %+v", b)
		c := a.MergeRight(iterable.FromSlice(b))
		t.Logf("c = %+v", c)
		require.NoError(t, c.checkInvariant())

		for _, e := range b {
			require.Equal(t, e.Value, c.GetOrZero(e.Key), "left value for key %+v", e.Key)
		}

		for it := iterable.Start[dict.Entry[key, int]](a); it.HasNext(); it.Next() {
			e := it.Current()
			if !a.ContainsKey(e.Key) {
				require.Equal(t, e.Value, c.GetOrZero(e.Key), "right value for key %+v", e.Key)
			}
		}
	})
}

func TestMergeLeft2(t *testing.T) {
	a := New(keyHash, dict.Entry[key, int]{Key: "", Value: 0}, dict.Entry[key, int]{Key: "b", Value: 0})
	b := New(keyHash, dict.Entry[key, int]{Key: "a", Value: 0}, dict.Entry[key, int]{Key: "aa", Value: 0})
	c := a.MergeLeft(b)
	t.Logf("a = %+v", a)
	t.Logf("b = %+v", b)
	t.Logf("c = %+v", c)
	t.Logf("a.root = %+v", a.root)
	t.Logf("b.root = %+v", b.root)
	t.Logf("c.root = %+v", c.root)
	expected := New(keyHash, dict.Entry[key, int]{Key: "", Value: 0}, dict.Entry[key, int]{Key: "b", Value: 0},
		dict.Entry[key, int]{Key: "a", Value: 0}, dict.Entry[key, int]{Key: "aa", Value: 0})
	require.Equal(t, expected.String(), c.String())
	//require.True(t, expected.Equals(c))
}

func TestFilterMap(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		a := genDict().Draw(t, "").(Dict[key, int])
		b := a.FilterMap(func(_ key, v int) (int, bool) {
			if v%2 == 0 {
				return 10 + v, true
			}
			return 0, false
		})
		t.Logf("a = %+v", a)
		t.Logf("b = %+v", b)
		t.Logf("a.root = %+v", a.root)
		t.Logf("b.root = %+v", b.root)
		require.NoError(t, a.checkInvariant())
		require.NoError(t, b.checkInvariant())

		for it := iterable.Start[dict.Entry[key, int]](a); it.HasNext(); it.Next() {
			if it.Current().Value%2 == 0 {
				require.True(t, b.ContainsKey(it.Current().Key), "key %+v does not exist in result", it.Current().Key)
				require.Equal(t, it.Current().Value+10, b.GetOrZero(it.Current().Key))
			}
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

func TestMerge2(t *testing.T) {
	dict0 := New[key, int](keyHash)
	dict1 := dict0.Set("", 0)
	dict2 := dict0.Set("a", 0)
	dict3 := dict1.Set("aa", 0)
	dict4 := dict2.MergeLeft(dict3)
	log.Printf("dict4 = %+v", dict4)
	require.True(t, dict4.ContainsKey("a"))
}

func TestMerge3(t *testing.T) {
	a := New[key, int](keyHash)
	b := New[key, int](keyHash, dict.Entry[key, int]{"abbb", 1}, dict.Entry[key, int]{"ababbb", 2})
	c := a.MergeLeft(b)
	log.Printf("a = %+v", a)
	log.Printf("b = %+v", b)
	log.Printf("c = %+v", c)
	log.Printf("a = %+v", a.root)
	log.Printf("b = %+v", b.root)
	log.Printf("c = %+v", c.root)
	require.NoError(t, a.checkInvariant(), "a invariant")
	require.NoError(t, b.checkInvariant(), "b invariant")
	require.NoError(t, c.checkInvariant(), "c invariant")
}
