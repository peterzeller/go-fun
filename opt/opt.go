package opt

import (
	"encoding/json"
	"fmt"

	"github.com/peterzeller/go-fun/iterable"
)

// Optional is a type that holds an optional value of type T.
// Optional values are constructed with either the None or the Some function.
type Optional[T any] struct {
	present bool
	value   T
}

// None returns an Optional without a value
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// Some returns an Optional with a present value
func Some[T any](v T) Optional[T] {
	return Optional[T]{
		present: true,
		value:   v,
	}
}

// First returns the first value from the given iterable or None if the iterable is empty
func First[T any](it iterable.Iterable[T]) Optional[T] {
	i := it.Iterator()
	if v, ok := i.Next(); ok {
		return Some(v)
	}
	return None[T]()
}

// Present returns true if this optional has a value
func (o Optional[T]) Present() bool {
	return o.present
}

// Get returns the value stored in this optional, and a boolean that is true only if the value is present.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// GetPointer returns the value stored in this optional as a pointer, using nil to represent the absent value
func (o Optional[T]) GetPointer() *T {
	if o.present {
		v := o.value
		return &v
	}
	return nil
}

// Value returns the value stored in this optional.
// Returns the zero value if no value is present
func (o Optional[T]) Value() T {
	return o.value
}

// OrElse returns the value stored if present, or else the given default value
func (o Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// OrElseGet returns the value stored if present, or else the value returned by the given default value function
func (o Optional[T]) OrElseGet(defaultValue func() T) T {
	if o.present {
		return o.value
	}
	return defaultValue()
}

// OrElsePanic returns the value stored if present, or else panics
func (o Optional[T]) OrElsePanic() T {
	if o.present {
		return o.value
	}
	panic(ErrNoValue)
}

var ErrNoValue = fmt.Errorf("getting value from empty optional")

// Iterator to iterate over the elements in this Optional value.
func (o Optional[T]) Iterator() iterable.Iterator[T] {
	if o.present {
		return iterable.Singleton(o.value).Iterator()
	}
	return iterable.Empty[T]().Iterator()
}

// String representation of the optional value
func (o Optional[T]) String() string {
	if o.present {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None()"
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (o *Optional[T]) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		*o = None[T]()
		return nil
	}
	o.present = true
	return json.Unmarshal(bytes, &o.value)
}

// MarshalJSON implements the json.Marshaler interface
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.present {
		return json.Marshal(o.value)
	}
	return []byte(`null`), nil
}
