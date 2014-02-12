package seq

import (
	. "testing"
)

// Test creating a list and calling the Seq interface methods on it
func TestListSeq(t *T) {
	ints := []interface{}{1, "a", 5.0}
	intsl := uint64(len(ints))

	// Testing creation and Size
	l := NewList(ints...)

	// Testing FirstRest, Size, and ToSlice
	sl := Seq(l)
	for i := range ints {
		assertValue(sl.Size(), intsl-uint64(i), t)
		assertSeqContents(sl, ints[i:], t)

		first, rest, ok := sl.FirstRest()
		assertValue(ok, true, t)
		assertValue(first, ints[i], t)

		sl = rest
	}

	// sl should be empty at this point. We use nilpointer because checking a
	// nil pointer against nil after both have been wrapped in an interface{}
	// (as happens when passed to assertValue) causes equality to not work.
	l = sl.ToList()
	var nilpointer *List
	assertEmpty(l, t)
	assertValue(l.el, nil, t)
	assertValue(l.next, nilpointer, t)
	assertValue(len(l.ToSlice()), 0, t)

	// Testing creation of empty List. We dereference the pointers so that we're
	// testing the actual values inside the List structs.
	emptyl := NewList()
	assertValue(*emptyl, *l, t)
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
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{4, 3, 2, 1, 0}, t)

	// Degenerate case
	l = NewList()
	nl = l.Prepend(0)
	assertEmpty(l, t)
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
	assertSeqContents(nl, intl1, t)

	nl = l1.PrependSeq(blank1)
	assertEmpty(blank1, t)
	assertSeqContents(nl, intl1, t)
}

// Test appending to the end of a List
func TestAppend(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1}
	l := NewList(intl...)
	nl := l.Append(0)
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{3, 2, 1, 0}, t)

	// Edge case (algorithm gets weird here)
	l = NewList(1)
	nl = l.Append(0)
	assertSeqContents(l, []interface{}{1}, t)
	assertSeqContents(nl, []interface{}{1, 0}, t)

	// Degenerate case
	l = NewList()
	nl = l.Append(0)
	assertEmpty(l, t)
	assertSeqContents(nl, []interface{}{0}, t)
}

// Test reversing a List
func TestReverse(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1}
	l := NewList(intl...)
	nl := l.Reverse()
	assertSeqContents(l, intl, t)
	assertSeqContents(nl, []interface{}{1, 2, 3}, t)

	// Degenerate case
	l = NewList()
	nl = l.Reverse()
	assertEmpty(l, t)
	assertEmpty(nl, t)
}
