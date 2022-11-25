package promise

import "context"

// SetupMocks replaces the functions in this file with mocks or fake functions.
// it returns a function that can be called to reset the functions to their default value.
func SetupMocks(startGoroutineFunc func(f func()), waitForContextsFunc func(ctxA context.Context, ctxB context.Context) bool) (cancel func()) {
	startGoroutine = startGoroutineFunc
	waitForContexts = waitForContextsFunc
	return func() {
		startGoroutine = startGoroutineDefault
		waitForContexts = waitForContextsDefault
	}
}

// startGoroutineDefault starts a new goroutine
func startGoroutineDefault(f func()) {
	go f()
}

// waitForContextsDefault for either context A or context B to finish.
// Returns true if ctxA finishes first, false otherwise.
func waitForContextsDefault(ctxA, ctxB context.Context) bool {
	select {
	case <-ctxA.Done():
		return true
	case <-ctxB.Done():
		return false
	}
}

var startGoroutine = startGoroutineDefault
var waitForContexts = waitForContextsDefault
