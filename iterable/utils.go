package iterable

import (
	"fmt"
	"strings"
)

func String[T any](i Iterable[T]) string {
	var res strings.Builder
	res.WriteString("[")
	first := true
	for it := Start(i); it.HasNext(); it.Next() {
		if !first {
			res.WriteString(", ")
		}
		res.WriteString(fmt.Sprintf("%+v", it.current))
		first = false
	}
	res.WriteString("]")
	return res.String()
}
