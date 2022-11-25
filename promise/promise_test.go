package promise_test

import (
	"context"
	"fmt"
	"github.com/peterzeller/go-fun/promise"
)

func ExamplePromise_Resolve() {
	p := promise.New[int]()
	go func() {
		p.Resolve(42)
	}()
	v, err := p.Future().Wait(context.Background())
	fmt.Printf("v = %v, err = %v\n", v, err)
	// output: v = 42, err = <nil>
}

func ExamplePromise_Reject() {
	p := promise.New[int]()
	reason := fmt.Errorf("example rejection")
	go func() {
		p.Reject(reason)
	}()
	v, err := p.Future().Wait(context.Background())
	fmt.Printf("v = %v, err = %v\n", v, err)
	// output: v = 0, err = example rejection
}
