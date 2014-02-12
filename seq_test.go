package seq

import (
	. "testing"
)

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
