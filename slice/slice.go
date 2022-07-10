package slice

import "github.com/peterzeller/go-fun/equality"

// ContainsEq checks whether a slice contains an element
func ContainsEq[T any](s []T, elem T, eq equality.Equality[T]) bool {
	for _, e := range s {
		if eq.Equal(e, elem) {
			return true
		}
	}
	return false
}

// Contains checks whether a slice contains an element
func Contains[T comparable](s []T, elem T) bool {
	for _, e := range s {
		if e == elem {
			return true
		}
	}
	return false
}

// Exists checks whether some element in the slice fulfills the given condition
func Exists[T any](s []T, cond func(T) bool) bool {
	for _, e := range s {
		if cond(e) {
			return true
		}
	}
	return false
}

// Forall checks whether all elements in the slice fulfill the given condition
func Forall[T any](s []T, cond func(T) bool) bool {
	for _, e := range s {
		if !cond(e) {
			return false
		}
	}
	return true
}

// Equal checks whether two slices are equal
func Equal[T any](a []T, b []T, eq equality.Equality[T]) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !eq.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

// PrefixOf checks whether a is a prefix of b
func PrefixOf[T any](a, b []T, eq equality.Equality[T]) bool {
	if len(a) > len(b) {
		return false
	}
	for i := range a {
		if !eq.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

// IndexOf returns the index of the first occurrence of elem in the slice elems or -1 if elem is not in the slice.
func IndexOf[T any](elem T, elems []T, eq equality.Equality[T]) int {
	for i, t := range elems {
		if eq.Equal(t, elem) {
			return i
		}
	}
	return -1
}

// RemoveAt removes the element at the given index from the slice and returns the modified slice.
func RemoveAt[T any](s []T, index int) []T {
	return append(append([]T{}, s[:index]...), s[index+1:]...)
}

// RemoveFirst removes the first occurrence of an element from the slice and returns the modified slice.
func RemoveFirst[T any](s []T, elem T, eq equality.Equality[T]) []T {
	res := make([]T, 0, len(s))
	removed := false
	for _, t := range s {
		if !removed && eq.Equal(t, elem) {
			removed = true
		} else {
			res = append(res, t)
		}
	}
	return res
}

// RemoveAll removes all occurrences of the element from the slice and returns the modified slice.
func RemoveAll[T any](s []T, elem T, eq equality.Equality[T]) []T {
	res := make([]T, 0, len(s))
	for _, t := range s {
		if !eq.Equal(t, elem) {
			res = append(res, t)
		}
	}
	return res
}
