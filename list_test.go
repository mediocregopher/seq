package seq

import (
	. "testing"

	"github.com/stretchr/testify/assert"
)

// Asserts that the given list is properly formed and has all of its size fields
// filled in correctly
func assertSaneList(l *List, t *T) {
	if Size(l) == 0 {
		var nilpointer *List
		assert.Equal(t, nilpointer, l)
		return
	}

	size := Size(l)
	assert.Equal(t, size-1, Size(l.next))
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
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, nilpointer, l)
	assert.Equal(t, 0, len(ToSlice(l)))

	// Testing creation of empty List.
	emptyl := NewList()
	assert.Equal(t, nilpointer, emptyl)
}

// Test the string representation of a List
func TestStringSeq(t *T) {
	l := NewList(0, 1, 2, 3)
	assert.Equal(t, "( 0 1 2 3 )", l.String())

	l = NewList(0, 1, 2, NewList(3, 4), 5, NewList(6, 7, 8))
	assert.Equal(t, "( 0 1 2 ( 3 4 ) 5 ( 6 7 8 ) )", l.String())
}

// Test prepending an element to the beginning of a list
func TestPrepend(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1, 0}
	l := NewList(intl...)
	nl := l.Prepend(4)
	assertSaneList(l, t)
	assertSaneList(nl, t)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{4, 3, 2, 1, 0}, ToSlice(nl))

	// Degenerate case
	l = NewList()
	nl = l.Prepend(0)
	assert.Equal(t, 0, Size(l))
	assertSaneList(nl, t)
	assert.Equal(t, []interface{}{0}, ToSlice(nl))
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
	assert.Equal(t, intl1, ToSlice(l1))
	assert.Equal(t, intl2, ToSlice(l2))
	assert.Equal(t, []interface{}{0, 1, 2, 3, 4}, ToSlice(nl))

	// Degenerate cases
	blank1 := NewList()
	blank2 := NewList()
	nl = blank1.PrependSeq(blank2)
	assert.Equal(t, 0, Size(blank1))
	assert.Equal(t, 0, Size(blank2))
	assert.Equal(t, 0, Size(nl))

	nl = blank1.PrependSeq(l1)
	assert.Equal(t, 0, Size(blank1))
	assertSaneList(nl, t)
	assert.Equal(t, intl1, ToSlice(nl))

	nl = l1.PrependSeq(blank1)
	assert.Equal(t, 0, Size(blank1))
	assertSaneList(nl, t)
	assert.Equal(t, intl1, ToSlice(nl))
}

// Test appending to the end of a List
func TestAppend(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1}
	l := NewList(intl...)
	nl := l.Append(0)
	assertSaneList(l, t)
	assertSaneList(nl, t)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{3, 2, 1, 0}, ToSlice(nl))

	// Edge case (algorithm gets weird here)
	l = NewList(1)
	nl = l.Append(0)
	assertSaneList(l, t)
	assertSaneList(nl, t)
	assert.Equal(t, []interface{}{1}, ToSlice(l))
	assert.Equal(t, []interface{}{1, 0}, ToSlice(nl))

	// Degenerate case
	l = NewList()
	nl = l.Append(0)
	assert.Equal(t, 0, Size(l))
	assertSaneList(nl, t)
	assert.Equal(t, []interface{}{0}, ToSlice(nl))
}

// Test retrieving items from a List
func TestNth(t *T) {
	// Normal case, in bounds
	intl := []interface{}{0, 2, 4, 6, 8}
	l := NewList(intl...)
	r, ok := l.Nth(3)
	assertSaneList(l, t)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 6, r)
	assert.Equal(t, true, ok)

	// Normal case, out of bounds
	r, ok = l.Nth(8)
	assertSaneList(l, t)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, nil, r)
	assert.Equal(t, false, ok)

	// Degenerate case
	l = NewList()
	r, ok = l.Nth(0)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, nil, r)
	assert.Equal(t, false, ok)
}

// Test that two lists compare equality correctly
func TestListEqual(t *T) {
	// Degenerate case
	l1, l2 := NewList(), NewList()
	assert.Equal(t, true, l1.Equal(l2))
	assert.Equal(t, true, l2.Equal(l1))

	// False with different sizes
	l1 = l1.Prepend(1)
	assert.Equal(t, false, l1.Equal(l2))
	assert.Equal(t, false, l2.Equal(l1))

	// False with same sizes
	l2 = l2.Prepend(2)
	assert.Equal(t, false, l1.Equal(l2))
	assert.Equal(t, false, l2.Equal(l1))

	// Now true
	l1 = l1.Prepend(2)
	l2 = l2.Append(1)
	assert.Equal(t, true, l1.Equal(l2))
	assert.Equal(t, true, l2.Equal(l1))

	// False with embedded list
	l1 = l1.Prepend(NewList(3))
	assert.Equal(t, false, l1.Equal(l2))
	assert.Equal(t, false, l2.Equal(l1))

	// True with embedded set
	l2 = l2.Prepend(NewList(3))
	assert.Equal(t, true, l1.Equal(l2))
	assert.Equal(t, true, l2.Equal(l1))
}
