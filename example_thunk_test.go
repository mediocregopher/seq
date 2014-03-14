package seq

import (
	"fmt"
	"github.com/mediocregopher/seq"
)

// numberThunk is used to create a sequence of inifinte, sequential integers
func numberThunk(i int) seq.Thunk {
	return func() (interface{}, seq.Thunk, bool) {
		return i, numberThunk(i + 1), true
	}
}

// Numbers is a nice wrapper around numberThunk which creates an root thunk and
// wraps it with a Lazy
func Numbers() *seq.Lazy {
	rootThunk := numberThunk(0)
	return seq.NewLazy(rootThunk)
}

func ExampleThunk() {
	var el interface{}
	var s seq.Seq = Numbers()
	var ok bool
	for i := 0; i < 10; i++ {
		if el, s, ok = s.FirstRest(); ok {
			fmt.Println(el)
		} else {
			break
		}
	}
}
