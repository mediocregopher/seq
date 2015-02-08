package main

// Trivial example of building a list and then mapping over it, then filtering
// it.

import (
	"fmt"

	"github.com/mediocregopher/seq"
)

func Inc(i interface{}) interface{} {
	return i.(int) + 1
}

func Even(i interface{}) bool {
	return i.(int)%2 == 0
}

func main() {

	// Some unnecessary list construction
	var l seq.Seq = seq.NewList(0, 1, 2, 3)
	lpost := seq.NewList()
	for i := 4; i < 10; i++ {
		lpost = lpost.Prepend(i)
	}
	l = seq.ToList(seq.Reverse(lpost)).PrependSeq(l)
	// At this point l = (0 1 2 3 4 5 6 8 9)

	lMF := seq.Filter(Even, seq.Map(Inc, l))
	lFM := seq.Map(Inc, seq.Filter(Even, l))

	// Print the original list and all the results
	fmt.Println(l)
	fmt.Println(lMF)
	fmt.Println(lFM)

}
