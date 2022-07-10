package hashset_test

import (
	"fmt"

	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/set/hashset"
)

func ExampleSet_Contains() {
	s := hashset.New(hash.String(), "a", "b", "c")
	fmt.Printf("contains a: %v\n", s.Contains("a"))
	fmt.Printf("contains x: %v\n", s.Contains("x"))
	// output:
	// contains a: true
	// contains x: false
}

func ExampleSet_Add() {
	s1 := hashset.New(hash.String(), "a", "b", "c")
	s2 := s1.Add("x", "y")
	fmt.Printf("s1 = %v\n", s1)
	fmt.Printf("s2 = %v\n", s2)
	// output:
	// s1 = [b, a, c]
	// s2 = [b, x, a, c, y]
}

func ExampleSet_Remove() {
	s1 := hashset.New(hash.String(), "a", "b", "c")
	s2 := s1.Remove("c")
	fmt.Printf("s1 = %v\n", s1)
	fmt.Printf("s2 = %v\n", s2)
	// output:
	// s1 = [b, a, c]
	// s2 = [b, a]
}

func ExampleSet_Iterator() {
	s := hashset.New(hash.Num[int](), 1, 2, 3)
	it := s.Iterator()
	sum := 0
	for {
		v, ok := it.Next()
		if !ok {
			break
		}
		sum += v
	}
	fmt.Printf("sum = %v\n", sum)
	// output: sum = 6
}

func ExampleSet_Union() {
	s1 := hashset.New(hash.String(), "a", "b", "c")
	s2 := hashset.New(hash.String(), "b", "d")
	s3 := s1.Union(s2)
	fmt.Printf("s1 = %v\n", s1)
	fmt.Printf("s2 = %v\n", s2)
	fmt.Printf("s3 = %v\n", s3)
	// output:
	// s1 = [b, a, c]
	// s2 = [b, d]
	// s3 = [b, a, c, d]
}

func ExampleSet_Intersect() {
	s1 := hashset.New(hash.String(), "a", "b", "c")
	s2 := hashset.New(hash.String(), "b", "c", "d")
	s3 := s1.Intersect(s2)
	fmt.Printf("s1 = %v\n", s1)
	fmt.Printf("s2 = %v\n", s2)
	fmt.Printf("s3 = %v\n", s3)
	// output:
	// s1 = [b, a, c]
	// s2 = [b, c, d]
	// s3 = [b, c]
}

func ExampleSet_Minus() {
	s1 := hashset.New(hash.String(), "a", "b", "c", "d")
	s2 := hashset.New(hash.String(), "b", "d", "e")
	s3 := s1.Minus(s2)
	fmt.Printf("s1 = %v\n", s1)
	fmt.Printf("s2 = %v\n", s2)
	fmt.Printf("s3 = %v\n", s3)
	// output:
	// s1 = [b, a, c, d]
	// s2 = [e, b, d]
	// s3 = [a, c]
}
