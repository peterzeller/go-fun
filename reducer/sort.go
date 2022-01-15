package reducer

import (
	"github.com/peterzeller/go-fun/mutable"
	"github.com/peterzeller/go-fun/zero"
)

func Sorted[A, B any](cmp func(A, A) bool, next Reducer[A, B]) Reducer[A, B] {
	return func() ReducerInstance[A, B] {
		inputs := make([]A, 0)
		nextI := next()
		return ReducerInstance[A, B]{
			Complete: func() B {
				qs := quicksortInit(inputs, cmp)
				for {
					a, ok := qs.Next()
					if !ok {
						return nextI.Complete()
					}
					cont := nextI.Step(a)
					if !cont {
						return nextI.Complete()
					}
				}
			},
			Step: func(a A) bool {
				inputs = append(inputs, a)
				return true
			},
		}
	}
}

type quicksortTodo[A any] struct {
	// stack with the leftmost slice on the top
	stack *mutable.Stack[[]A]
	cmp   func(A, A) bool
}

func quicksortInit[A any](slice []A, cmp func(A, A) bool) *quicksortTodo[A] {
	return &quicksortTodo[A]{
		stack: mutable.NewStack(slice),
		cmp:   cmp,
	}

}

func (q *quicksortTodo[A]) Next() (A, bool) {
	for !q.stack.Empty() {
		slice := q.stack.Pop()
		if len(slice) == 0 {
			continue
		}
		if len(slice) == 1 {
			return slice[0], true
		}
		first, mid, second := quicksortPartition(slice, q.cmp)
		q.stack.Push(second)
		q.stack.Push(mid)
		q.stack.Push(first)
	}
	return zero.Value[A](), false
}

func quicksortPartition[A any](ar []A, cmp func(A, A) bool) ([]A, []A, []A) {
	if len(ar) == 0 {
		return []A{}, []A{}, []A{}
	}
	if len(ar) > 3 {
		m := len(ar) / 2
		hi := len(ar)
		if len(ar) > 40 {
			s := len(ar) / 8
			medianOfThree(cmp, &ar[0], &ar[s], &ar[2*s])
			medianOfThree(cmp, &ar[m], &ar[m-s], &ar[m+s])
			medianOfThree(cmp, &ar[hi-1], &ar[hi-1-s], &ar[hi-1-2*s])
		}
		medianOfThree(cmp, &ar[0], &ar[m], &ar[hi-1])
	}

	pivot := ar[0]
	pivotPos := 0

	for i := 1; i < len(ar); i++ {
		if cmp(ar[i], pivot) {
			if i == pivotPos+1 {
				ar[pivotPos], ar[i] = ar[i], ar[pivotPos]
			} else {
				ar[pivotPos], ar[pivotPos+1], ar[i] = ar[i], ar[pivotPos], ar[pivotPos+1]
			}
			pivotPos++
		}
	}
	return ar[0:pivotPos], []A{pivot}, ar[pivotPos+1:]
}

func medianOfThree[A any](less func(A, A) bool, m1, m0, m2 *A) {
	if less(*m1, *m0) {
		*m1, *m0 = *m0, *m1
	}
	if less(*m2, *m1) {
		*m2, *m1 = *m1, *m2
		if less(*m1, *m0) {
			*m1, *m0 = *m0, *m1
		}
	}
}
