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
