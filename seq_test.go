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
		t.Fatalf("New slice wrong: %v", nl.ToSlice())
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

// Test the Any function
func TestAny(t *T) {
	// Normal cases
	int1 := []interface{}{1, 2, 3, 4}
	int2 := []interface{}{1, 2, 3}
	l1 := NewList(int1...)
	l2 := NewList(int2...)
	fn := func(el interface{}) bool {
		return el.(int) > 3
	}

	r, ok := Any(fn, l1)
	if !intSlicesEq(l1.ToSlice(), int1) {
		t.Fatalf("Original slice changed: %v", l1.ToSlice())
	}
	if r != 4 || !ok {
		t.Fatalf("Returned values wrong: %v, %v", r, ok)
	}

	r, ok = Any(fn, l2)
	if !intSlicesEq(l2.ToSlice(), int2) {
		t.Fatalf("Original slice changed: %v", l2.ToSlice())
	}
	if r != nil || ok {
		t.Fatalf("Returned values wrong: %v, %v", r, ok)
	}

	// Degenerate cases
	l := NewList()
	r, ok = Any(fn, l)
	if !intSlicesEq(l.ToSlice(), []interface{}{}) {
		t.Fatalf("Degenerate slice changed: %v", l.ToSlice())
	}
	if r != nil || ok {
		t.Fatalf("Degenerate values wrong: %v, %v", r, ok)
	}
}

// Test the Filter function
func TestFilter(t *T) {
	// Normal case
	intl := []interface{}{1, 2, 3, 4, 5}
	l := NewList(intl...)
	fn := func(el interface{}) bool {
		return el.(int) % 2 != 0
	}

	r := Filter(fn, l)
	if !intSlicesEq(l.ToSlice(), intl) {
		t.Fatalf("Original slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(r.ToSlice(), []interface{}{1, 3, 5}) {
		t.Fatalf("Returned slice wrong: %v", r.ToSlice())
	}

	// Degenerate cases
	l = NewList()
	r = Filter(fn, l)
	if !intSlicesEq(l.ToSlice(), []interface{}{}) {
		t.Fatalf("Degenerate slice changed: %v", l.ToSlice())
	}
	if len(r.ToSlice()) != 0 {
		t.Fatalf("Degenerate return wrong: %v", r.ToSlice())
	}
}
