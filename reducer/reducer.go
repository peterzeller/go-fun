package reducer

import "github.com/peterzeller/go-fun/v2/iterable"

// Reducer works on streams of A and produces a B
type Reducer[A, B any] func() ReducerInstance[A, B]

// ReducerInstance is a temporary data structure created by the Reducer at the beginning of a pipeline.
type ReducerInstance[A, B any] struct {
	Complete func() B
	Step     func(A) bool
}

func (r Reducer[A, B]) Apply(i iterable.Iterable[A]) B {
	return r.ApplyIterator(i.Iterator())
}

func (r Reducer[A, B]) ApplyIterator(it iterable.Iterator[A]) B {
	ri := r()
	for {
		a, ok := it.Next()
		if !ok {
			break
		}
		cont := ri.Step(a)
		if !cont {
			break
		}
	}
	return ri.Complete()
}

func (r Reducer[A, B]) ApplySlice(as []A) B {
	ri := r()
	for _, a := range as {
		cont := ri.Step(a)
		if !cont {
			break
		}
	}
	return ri.Complete()
}
