package seq

import (
	. "testing"
)

// Asserts that the given list is properly formed and has all of its size fields
// filled in correctly
func assertSaneList(l *List, t *T) {
	if Size(l) == 0 {
		var nilpointer *List
		assertValue(l, nilpointer, t)
		return
	}

	size := Size(l)
	assertValue(Size(l.next), size - 1, t)
	assertSaneList(l.next, t)
}

// Test creating a list and calling the Seq interface methods on it
func TestListSeq(t *T) {
	ints := []interface{}{1, "a", 5.0}

	// Testing creation and Seq interface methods
	l := NewList(ints...)
	sl := testSeqGen(t, l, ints)

	// sl should be empty at this point
	l = ToList(sl)
	var nilpointer *List
	assertEmpty(l, t)
	assertValue(l, nilpointer, t)
	assertValue(len(ToSlice(l)), 0, t)

	// Testing creation of empty List.
	emptyl := NewList()
	assertValue(emptyl, nilpointer, t)
}

// Test the string representation of a List
func TestStringSeq(t *T) {
	l := NewList(0, 1, 2, 3)
	assertValue(l.String(), "( 0 1 2 3 )", t)

	l = NewList(0, 1, 2, NewList(3, 4), 5, NewList(6, 7, 8))
	assertValue(l.String(), "( 0 1 2 ( 3 4 ) 5 ( 6 7 8 ) )", t)
}

// Test prepending an element to the beginning of a list
func TestPrepend(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1, 0}
	l := NewList(intl...)
	nl := l.Prepend(4)
	assertSaneList(l, t)
	assertSaneList(nl, t)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{4, 3, 2, 1, 0}, t)

	// Degenerate case
	l = NewList()
	nl = l.Prepend(0)
	assertEmpty(l, t)
	assertSaneList(nl, t)
	assertSeqContents(nl, []interface{}{0}, t)
}

// Test prepending a Seq to the beginning of a list
func TestPrependSeq(t *T) {
	//Normal case
	intl1 := []interface{}{3, 4}
	intl2 := []interface{}{0, 1, 2}
	l1 := NewList(intl1...)
	l2 := NewList(intl2...)
	nl := l1.PrependSeq(l2)
	assertSaneList(l1, t)
	assertSaneList(l2, t)
	assertSaneList(nl, t)
	assertSeqContents(l1, intl1, t)
	assertSeqContents(l2, intl2, t)
	assertSeqContents(nl, []interface{}{0, 1, 2, 3, 4}, t)

	// Degenerate cases
	blank1 := NewList()
	blank2 := NewList()
	nl = blank1.PrependSeq(blank2)
	assertEmpty(blank1, t)
	assertEmpty(blank2, t)
	assertEmpty(nl, t)

	nl = blank1.PrependSeq(l1)
	assertEmpty(blank1, t)
	assertSaneList(nl, t)
	assertSeqContents(nl, intl1, t)

	nl = l1.PrependSeq(blank1)
	assertEmpty(blank1, t)
	assertSaneList(nl, t)
	assertSeqContents(nl, intl1, t)
}

// Test appending to the end of a List
func TestAppend(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1}
	l := NewList(intl...)
	nl := l.Append(0)
	assertSaneList(l, t)
	assertSaneList(nl, t)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{3, 2, 1, 0}, t)

	// Edge case (algorithm gets weird here)
	l = NewList(1)
	nl = l.Append(0)
	assertSaneList(l, t)
	assertSaneList(nl, t)
	assertSeqContents(l, []interface{}{1}, t)
	assertSeqContents(nl, []interface{}{1, 0}, t)

	// Degenerate case
	l = NewList()
	nl = l.Append(0)
	assertEmpty(l, t)
	assertSaneList(nl, t)
	assertSeqContents(nl, []interface{}{0}, t)
}

// Test retrieving items from a List
func TestNth(t *T) {
	// Normal case, in bounds
	intl := []interface{}{0, 2, 4, 6, 8}
	l := NewList(intl...)
	r, ok := l.Nth(3)
	assertSaneList(l, t)
	assertSeqContents(l, intl, t)
	assertValue(r, 6, t)
	assertValue(ok, true, t)

	// Normal case, out of bounds
	r, ok = l.Nth(8)
	assertSaneList(l, t)
	assertSeqContents(l, intl, t)
	assertValue(r, nil, t)
	assertValue(ok, false, t)

	// Degenerate case
	l = NewList()
	r, ok = l.Nth(0)
	assertEmpty(l, t)
	assertValue(r, nil, t)
	assertValue(ok, false, t)
}
