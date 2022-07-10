package hash_test

import (
	"fmt"

	"github.com/peterzeller/go-fun/hash"
)

func ExampleNum() {
	h := hash.Num[int]()
	fmt.Println("hash of 42 =", h.Hash(42))
	// output: hash of 42 = 42
}

func ExampleString() {
	h := hash.String()
	fmt.Println("1. hash is", h.Hash("hello world"))
	fmt.Println("2. hash is", h.Hash("hello worle"))
	fmt.Println("3. hash is", h.Hash("gello world"))
	// output:
	// 1. hash is 8618312879776256743
	// 2. hash is 8618311780264628532
	// 3. hash is 1592153891954394282
}

type exampleStruct struct {
	A int
	B string
}

func ExampleGob() {
	h := hash.Gob[exampleStruct]()
	v := exampleStruct{
		A: 42,
		B: "hello",
	}
	fmt.Println("hash is", h.Hash(v))
	// output: hash is -4727944318207215451
}

type pair struct {
	a int
	b string
}

func (p pair) Equal(other pair) bool {
	return p == other
}

func (p pair) Hash() int64 {
	return hash.CombineHashes(
		hash.Num[int]().Hash(p.a),
		hash.String().Hash(p.b))
}

func ExampleNatural() {
	h := hash.Natural[pair]()
	v := pair{
		a: 42,
		b: "hello",
	}
	fmt.Println("hash is", h.Hash(v))
	// output: hash is -6615550055289272862
}

func ExampleMap() {
	h := hash.Map(hash.String(), func(bs []byte) string {
		return string(bs)
	})
	fmt.Println("hash is", h.Hash([]byte("hello world")))
	// output: hash is 8618312879776256743
}

func ExamplePairHash() {
	ph := hash.PairHash(hash.Num[int](), hash.String())
	h := hash.Map(ph, func(v pair) hash.Pair[int, string] {
		return hash.Pair[int, string]{
			A: v.a,
			B: v.b,
		}
	})
	v := pair{
		a: 42,
		b: "hello",
	}
	fmt.Println("hash is", h.Hash(v))
	// output: hash is -6615550055289272862
}

func ExampleCombineHashes() {
	fmt.Println("combined hash is", hash.CombineHashes(1, 2, 42))
	// output: combined hash is 30856
}
