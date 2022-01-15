package reducer

import "github.com/peterzeller/go-fun/iterable"

func Apply[A, B any](s iterable.Iterable[A], reducer Reducer[A, B]) B {
	i := reducer()
	it := s.Iterator()
	for {
		a, ok := it.Next()
		if !ok {
			return i.Complete()
		}
		cont := i.Step(a)
		if !cont {
			return i.Complete()
		}
	}
}
