package seq

import (
	. "testing"
)

// Test mapping over a Seq
func TestMap(t *T) {
	// Normal case
	intl := []interface{}{1, 2, 3}
	intnl := []interface{}{2, 3, 4}
	l := NewList(intl...)
	fn := func(n interface{}) interface{} {
		return n.(int) + 1
	}
	nl := Map(fn, l)
	if !intSlicesEq(l.ToSlice(), intl) {
		t.Fatalf("Original slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), intnl) {
		t.Fatal("New slice wrong: %v", nl.ToSlice())
	}
	
	// Degenerate case
	l = NewList()
	nl = Map(fn, l)
	if !intSlicesEq(l.ToSlice(), []interface{}{}) {
		t.Fatalf("Degenerate slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), []interface{}{}) {
		t.Fatalf("New degenerate slice wrong: %v", nl.ToSlice())
	}
}

// Test reducing over a Seq
func TestReduce(t *T) {
	// Normal case
	intl := []interface{}{1, 2, 3, 4}
	l := NewList(intl...)
	fn := func(acc, el interface{}) (interface{}, bool) {
		return acc.(int) + el.(int), false
	}
	r := Reduce(fn, 0, l)
	if !intSlicesEq(l.ToSlice(), intl) {
		t.Fatalf("Original slice changed: %v", l.ToSlice())
	}
	if r != 10 {
		t.Fatalf("Reduced value wrong: %v", r)
	}

	// Short-circuit case
	fns := func(acc, el interface{}) (interface{}, bool) {
		return acc.(int) + el.(int), el.(int) > 2
	}
	r = Reduce(fns, 0, l)
	if !intSlicesEq(l.ToSlice(), intl) {
		t.Fatalf("(Short circuit) Original slice changed: %v", l.ToSlice())
	}
	if r != 6 {
		t.Fatalf("(Short circuit) Reduced value wrong: %v", r)
	}

	// Degenerate case
	l = NewList()
	r = Reduce(fn, 0, l)
	if !intSlicesEq(l.ToSlice(), []interface{}{}) {
		t.Fatalf("Degenerate original slice wrong: %v", l.ToSlice())
	}
	if r != 0 {
		t.Fatalf("Degenerate value wrong: %v", r)
	}
}
