package seq

import (
	. "testing"
)

// Test reversing a Seq
func TestReverse(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1}
	l := NewList(intl...)
	nl := Reverse(l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{1, 2, 3}, t)

	// Degenerate case
	l = NewList()
	nl = Reverse(l)
	assertEmpty(l, t)
	assertEmpty(nl, t)
}

// Test mapping over a Seq
func TestMap(t *T) {
	fn := func(n interface{}) interface{} {
		return n.(int) + 1
	}

	// Normal case
	intl := []interface{}{1, 2, 3}
	l := NewList(intl...)
	nl := Map(fn, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{2, 3, 4}, t)

	// Degenerate case
	l = NewList()
	nl = Map(fn, l)
	assertEmpty(l, t)
	assertEmpty(nl, t)
}

// Test reducing over a Seq
func TestReduce(t *T) {
	fn := func(acc, el interface{}) (interface{}, bool) {
		return acc.(int) + el.(int), false
	}

	// Normal case
	intl := []interface{}{1, 2, 3, 4}
	l := NewList(intl...)
	r := Reduce(fn, 0, l)
	assertSeqContents(l, intl, t)
	assertValue(r, 10, t)

	// Short-circuit case
	fns := func(acc, el interface{}) (interface{}, bool) {
		return acc.(int) + el.(int), el.(int) > 2
	}
	r = Reduce(fns, 0, l)
	assertSeqContents(l, intl, t)
	assertValue(r, 6, t)

	// Degenerate case
	l = NewList()
	r = Reduce(fn, 0, l)
	assertEmpty(l, t)
	assertValue(r, 0, t)
}

// Test the Any function
func TestAny(t *T) {
	fn := func(el interface{}) bool {
		return el.(int) > 3
	}

	// Value found case
	intl := []interface{}{1, 2, 3, 4}
	l := NewList(intl...)
	r, ok := Any(fn, l)
	assertSeqContents(l, intl, t)
	assertValue(r, 4, t)
	assertValue(ok, true, t)

	// Value not found case
	intl = []interface{}{1, 2, 3}
	l = NewList(intl...)
	r, ok = Any(fn, l)
	assertSeqContents(l, intl, t)
	assertValue(r, nil, t)
	assertValue(ok, false, t)

	// Degenerate case
	l = NewList()
	r, ok = Any(fn, l)
	assertEmpty(l, t)
	assertValue(r, nil, t)
	assertValue(ok, false, t)
}

// Test the Filter function
func TestFilter(t *T) {
	fn := func(el interface{}) bool {
		return el.(int)%2 != 0
	}

	// Normal case
	intl := []interface{}{1, 2, 3, 4, 5}
	l := NewList(intl...)
	r := Filter(fn, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(r, []interface{}{1, 3, 5}, t)

	// Degenerate cases
	l = NewList()
	r = Filter(fn, l)
	assertEmpty(l, t)
	assertEmpty(r, t)
}

// Test Flatten-ing of a Seq
func TestFlatten(t *T) {
	// Normal case
	intl1 := []interface{}{0, 1, 2}
	intl2 := []interface{}{3, 4, 5}
	l1 := NewList(intl1...)
	l2 := NewList(intl2...)
	blank := NewList()
	intl := []interface{}{-1, l1, l2, 6, blank, 7}
	l := NewList(intl...)
	nl := Flatten(l)
	assertSeqContents(l1, intl1, t)
	assertSeqContents(l2, intl2, t)
	assertEmpty(blank, t)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{-1, 0, 1, 2, 3, 4, 5, 6, 7}, t)

	// Degenerate case
	nl = Flatten(blank)
	assertEmpty(blank, t)
	assertEmpty(nl, t)
}

// Test taking from a Seq
func TestTake(t *T) {
	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4}
	l := NewList(intl...)
	nl := Take(3, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{0, 1, 2}, t)

	// Edge cases
	nl = Take(5, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, intl, t)

	nl = Take(6, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, intl, t)

	// Degenerate cases
	empty := NewList()
	nl = Take(1, empty)
	assertEmpty(empty, t)
	assertEmpty(nl, t)

	nl = Take(0, l)
	assertSeqContents(l, intl, t)
	assertEmpty(nl, t)
}

// Test taking from a Seq until a given condition
func TestTakeWhile(t *T) {
	pred := func(el interface{}) bool {
		return el.(int) < 3
	}

	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4, 5}
	l := NewList(intl...)
	nl := TakeWhile(pred, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{0, 1, 2}, t)

	// Edge cases
	intl = []interface{}{5, 5, 5}
	l = NewList(intl...)
	nl = TakeWhile(pred, l)
	assertSeqContents(l, intl, t)
	assertEmpty(nl, t)

	intl = []interface{}{0, 1, 2}
	l = NewList(intl...)
	nl = TakeWhile(pred, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{0, 1, 2}, t)

	// Degenerate case
	l = NewList()
	nl = TakeWhile(pred, l)
	assertEmpty(l, t)
	assertEmpty(nl, t)
}

// Test dropping from a Seq
func TestDrop(t *T) {
	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4}
	l := NewList(intl...)
	nl := Drop(3, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{3, 4}, t)

	// Edge cases
	nl = Drop(5, l)
	assertSeqContents(l, intl, t)
	assertEmpty(nl, t)

	nl = Drop(6, l)
	assertSeqContents(l, intl, t)
	assertEmpty(nl, t)

	// Degenerate cases
	empty := NewList()
	nl = Drop(1, empty)
	assertEmpty(empty, t)
	assertEmpty(nl, t)

	nl = Drop(0, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, intl, t)
}

// Test dropping from a Seq until a given condition
func TestDropWhile(t *T) {
	pred := func(el interface{}) bool {
		return el.(int) < 3
	}

	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4, 5}
	l := NewList(intl...)
	nl := DropWhile(pred, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{3, 4, 5}, t)

	// Edge cases
	intl = []interface{}{5, 5, 5}
	l = NewList(intl...)
	nl = DropWhile(pred, l)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, intl, t)

	intl = []interface{}{0, 1, 2}
	l = NewList(intl...)
	nl = DropWhile(pred, l)
	assertSeqContents(l, intl, t)
	assertEmpty(nl, t)

	// Degenerate case
	l = NewList()
	nl = DropWhile(pred, l)
	assertEmpty(l, t)
	assertEmpty(nl, t)
}
