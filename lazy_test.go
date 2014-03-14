package seq

import (
	. "testing"
	"time"
)

// Test lazy operation and thread-safety
func TestLazyBasic(t *T) {
	ch := make(chan int)
	mapfn := func(el interface{}) interface{} {
		i := el.(int)
		ch <- i
		return i
	}

	intl := []interface{}{0, 1, 2, 3, 4}
	l := NewList(intl...)
	ml := LMap(mapfn, l)

	for i := 0; i < 10; i++ {
		go func() {
			mlintl := ToSlice(ml)
			if !intSlicesEq(mlintl, intl) {
				panic("contents not right")
			}
		}()
	}

	for _, el := range intl {
		select {
		case elch := <-ch:
			assertValue(el, elch, t)
		case <-time.After(1 * time.Millisecond):
			t.Fatalf("Took too long reading result")
		}
	}
	close(ch)
}

// Test that arbitrary Seqs can turn into Lazy
func TestToLazy(t *T) {
	intl := []interface{}{0, 1, 2, 3, 4}
	l := NewList(intl...)
	ll := ToLazy(l)
	assertSeqContents(ll, intl, t)
}
