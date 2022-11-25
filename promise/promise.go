package promise

import (
	"context"
	"fmt"
	"github.com/peterzeller/go-fun/zero"
	"sync/atomic"
)

type Promise[T any] struct {
	data *pData[T]
}

type Future[T any] struct {
	data *pData[T]
}

type pData[T any] struct {
	// ctx is a context that is done when the promise is resolved
	ctx        context.Context
	cancelFunc func()
	err        atomic.Pointer[error]
	value      atomic.Pointer[T]
}

// New promise that can be resolved or rejected.
func New[T any]() Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return Promise[T]{data: &pData[T]{
		ctx:        ctx,
		cancelFunc: cancel,
	}}
}

// Async runs a function asynchronously and returns a future with the result returned by the function.
func Async[T any](f func() T) Future[T] {
	p := New[T]()
	startGoroutine(func() {
		defer func() {
			r := recover()
			if r != nil {
				switch err := r.(type) {
				case error:
					p.Reject(err)
				default:
					p.Reject(fmt.Errorf("%v", r))
				}
			}
		}()
		p.Resolve(f())
	})
	return p.Future()
}

// AsyncErr runs a function asynchronously and returns a future with the result or the error returned by the function.
func AsyncErr[T any](f func() (T, error)) Future[T] {
	p := New[T]()
	startGoroutine(func() {
		defer func() {
			r := recover()
			if r != nil {
				switch err := r.(type) {
				case error:
					p.Reject(err)
				default:
					p.Reject(fmt.Errorf("%v", r))
				}
			}
		}()
		t, err := f()
		if err == nil {
			p.Resolve(t)
		} else {
			p.Reject(err)
		}
	})
	return p.Future()
}

// AsyncVoid runs a function asynchronously and returns a future with no result
func AsyncVoid(f func()) Future[struct{}] {
	return Async[struct{}](func() struct{} {
		f()
		return struct{}{}
	})
}

func (p Promise[T]) Resolve(value T) {
	p.data.value.Store(&value)
	p.data.cancelFunc()
}

func (p Promise[T]) Reject(err error) {
	p.data.err.Store(&err)
	p.data.cancelFunc()
}

// Future bound to this promise
func (p Promise[T]) Future() Future[T] {
	return Future[T]{
		data: p.data,
	}
}

func (f Future[T]) Wait(ctx context.Context) (T, error) {
	if waitForContexts(f.data.ctx, ctx) {
		err := f.data.err.Load()
		if err != nil {
			return zero.Value[T](), *err
		}
		return *f.data.value.Load(), nil
	} else {
		return zero.Value[T](), ctx.Err()
	}
}
