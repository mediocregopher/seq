package seq

import (
	. "testing"
)

func intSlicesEq(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Test creating a list and calling the Seq interface methods on it
func TestListSeq(t *T) {
	ints := []interface{}{1, "a", 5.0}
	intsl := uint64(len(ints))

	// Testing creation and Size
	l := NewList(ints...)

	// Testing FirstRest, Size, and ToSlice
	sl := Seq(l)
	for i := range ints {
		if sl.Size() != intsl - uint64(i) {
			t.Fatalf("Size wrong: %v", sl.Size())
		} else if !intSlicesEq(sl.ToSlice(), ints[i:]) {
			t.Fatalf("ToSlice wrong: %v", sl.ToSlice())
		}

		first, rest, ok := sl.FirstRest();
		if !ok {
			t.Fatalf("ok wrong: %v", ok)
		} else if first != ints[i] {
			t.Fatalf("first wrong: %v", first)
		}
		sl = rest
	}

	// sl should be empty at this point
	l = sl.ToList()
	if l.Size() != 0 {
		t.Fatalf("Empty size wrong: %v", l.Size())
	} else if l.el != nil {
		t.Fatalf("Empty el wrong: %v", l.el)
	} else if l.next != nil {
		t.Fatalf("Empty next wrong: %v", l.next)
	} else if len(l.ToSlice()) != 0 {
		t.Fatalf("Empty toslice wrong: %v", l.ToSlice())
	}

	// Testing creation of empty l
	emptyl := NewList()
	if *emptyl != *l {
		t.Fatalf("Created empty wrong: %v", emptyl)
	}
}

// Test the string representation of a List
func TestStringSeq(t *T) {
	l := NewList(0, 1, 2, 3)
	if l.String() != "( 0 1 2 3 )" {
		t.Fatalf("String is wrong: %v", l.String())
	}

	nl := NewList(0, 1, 2, NewList(3, 4), 5, NewList(6, 7, 8))
	if nl.String() != "( 0 1 2 ( 3 4 ) 5 ( 6 7 8 ) )" {
		t.Fatalf("String is wrong: %v", nl.String())
	}
}

// Test prepending an element to the beginning of a list
func TestPrepend(t *T) {
	l := NewList(3, 2, 1, 0)
	nl := l.Prepend(4)
	if !intSlicesEq(l.ToSlice(), []interface{}{3, 2, 1, 0}) {
		t.Fatalf("Original slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), []interface{}{4, 3, 2, 1, 0}) {
		t.Fatalf("New slice wrong: %v", nl.ToSlice())
	}

	// Degenerate case
	l = NewList()
	nl = l.Prepend(0)
	if !intSlicesEq(nl.ToSlice(), []interface{}{0}) {
		t.Fatalf("New degenerate slice wrong: %v", nl.ToSlice())
	}
}

// Test prepending a Seq to the beginning of a list
func TestPrependSeq(t *T) {
	int1 := []interface{}{3, 4}
	int2 := []interface{}{0, 1, 2}
	l1 := NewList(int1...)
	l2 := NewList(int2...)
	nl := l1.PrependSeq(l2)

	if !intSlicesEq(l1.ToSlice(), int1) {
		t.Fatalf("First list changed: %v", l1.ToSlice())
	}
	if !intSlicesEq(l2.ToSlice(), int2) {
		t.Fatalf("Second list changed: %v", l2.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), []interface{}{0, 1, 2, 3, 4}) {
		t.Fatalf("New list wrong: %v", nl.ToSlice())
	}

	// Degenerate cases
	blank1 := NewList()
	blank2 := NewList()

	nl = blank1.PrependSeq(blank2)
	if !intSlicesEq(nl.ToSlice(), []interface{}{}) {
		t.Fatalf("Degenerate blank/blank case wrong: %v", nl.ToSlice())
	}

	nl = blank1.PrependSeq(l1)
	if !intSlicesEq(nl.ToSlice(), int1) {
		t.Fatalf("Degenerate blank/l1 case wrong: %v", nl.ToSlice())
	}

	nl = l1.PrependSeq(blank1)
	if !intSlicesEq(nl.ToSlice(), int1) {
		t.Fatalf("Degenerate l1/blank case wrong: %v", nl.ToSlice())
	}
}

// Test appending to the end of a List
func TestAppend(t *T) {
	// Normal case
	l := NewList(3, 2, 1)
	nl := l.Append(0)
	if !intSlicesEq(l.ToSlice(), []interface{}{3, 2, 1}) {
		t.Fatalf("Original slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), []interface{}{3, 2, 1, 0}) {
		t.Fatalf("New slice wrong: %v", nl.ToSlice())
	}

	// Edge case (algorithm gets weird here)
	l = NewList(1)
	nl = l.Append(0)
	if !intSlicesEq(l.ToSlice(), []interface{}{1}) {
		t.Fatalf("Original edge-case slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), []interface{}{1, 0}) {
		t.Fatalf("New edge-case slice wrong: %v", nl.ToSlice())
	}

	// Degenerate case
	l = NewList()
	nl = l.Append(0)
	if !intSlicesEq(nl.ToSlice(), []interface{}{0}) {
		t.Fatalf("New degenerate slice wrong: %v", nl.ToSlice())
	}
}

// Test reversing a List
func TestReverse(t *T) {
	// Normal case
	l := NewList(3, 2, 1)
	nl := l.Reverse()
	if !intSlicesEq(l.ToSlice(), []interface{}{3, 2, 1}) {
		t.Fatalf("Original slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), []interface{}{1, 2, 3}) {
		t.Fatalf("New slice wrong: %v", nl.ToSlice())
	}

	// Degenerate case
	l = NewList()
	nl = l.Reverse()
	if !intSlicesEq(l.ToSlice(), []interface{}{}) {
		t.Fatalf("Degenerate slice changed: %v", l.ToSlice())
	}
	if !intSlicesEq(nl.ToSlice(), []interface{}{}) {
		t.Fatalf("New degenerate slice wrong: %v", nl.ToSlice())
	}
}
