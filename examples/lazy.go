package main

// Trivial examples of using two lazy maps on a list of numbers to show both
// how they won't evaluate till necessary and how they can be chained together.

import (
	"fmt"
	"time"

	"github.com/mediocregopher/seq"
)

func F1(i interface{}) interface{} {
	fmt.Printf("multiplying %d by two\n", i.(int))
	return i.(int) * 2
}

func F2(i interface{}) interface{} {
	fmt.Printf("subtracting one from %d\n", i.(int))
	return i.(int) - 1
}

func main() {
	l := seq.NewList(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	ll := seq.LMap(F2, seq.LMap(F1, l))

	// Manually read items from the new mapped sequence, with a slight delay
	// between each read. You'll see that the print statements in F1/F2 are only
	// called as you're reading from the sequence, and that they are called
	// together, so no intermediate list is being created.
	forll := ll
	for {
		el, nextll, ok := forll.FirstRest()
		if !ok {
			break
		}
		forll = nextll
		fmt.Printf("el is %d\n", el)
		time.Sleep(500 * time.Millisecond)
	}

	// To show that the original list is unchanged and that the lazy list's
	// results are cached. F1/F2 aren't called again
	fmt.Println(l)
	fmt.Println(ll)
}
