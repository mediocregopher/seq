package seq

import (
	. "testing"

	"github.com/stretchr/testify/assert"
)

// Test creating a Set and calling the Seq interface methods on it
func TestSetSeq(t *T) {
	ints := []interface{}{1, "a", 5.0}

	// Testing creation and Seq interface methods
	s := NewSet(ints...)
	ss := testSeqNoOrderGen(t, s, ints)

	// ss should be empty at this point
	s = ToSet(ss)
	var nilpointer *Set
	assert.Equal(t, 0, Size(s))
	assert.Equal(t, nilpointer, s)
	assert.Equal(t, 0, len(ToSlice(s)))
}

// Test setting a value on a Set
func TestSetVal(t *T) {
	ints := []interface{}{0, 1, 2, 3, 4}
	ints1 := []interface{}{0, 1, 2, 3, 4, 5}

	// Degenerate case
	s := NewSet()
	assert.Equal(t, 0, Size(s))
	s, ok := s.SetVal(0)
	assertSeqContentsSet(t, []interface{}{0}, s)
	assert.Equal(t, true, ok)

	s = NewSet(ints...)
	s1, ok := s.SetVal(5)
	assertSeqContentsSet(t, ints, s)
	assertSeqContentsSet(t, ints1, s1)
	assert.Equal(t, true, ok)

	s2, ok := s1.SetVal(5)
	assertSeqContentsSet(t, ints1, s1)
	assertSeqContentsSet(t, ints1, s2)
	assert.Equal(t, false, ok)
}

// Test deleting a value from a Set
func TestDelVal(t *T) {
	ints := []interface{}{0, 1, 2, 3, 4}
	ints1 := []interface{}{0, 1, 2, 3}
	ints2 := []interface{}{1, 2, 3, 4}
	ints3 := []interface{}{1, 2, 3, 4, 5}

	// Degenerate case
	s := NewSet()
	assert.Equal(t, 0, Size(s))
	s, ok := s.DelVal(0)
	assert.Equal(t, 0, Size(s))
	assert.Equal(t, false, ok)

	s = NewSet(ints...)
	s1, ok := s.DelVal(4)
	assertSeqContentsSet(t, ints, s)
	assertSeqContentsSet(t, ints1, s1)
	assert.Equal(t, true, ok)

	s1, ok = s1.DelVal(4)
	assertSeqContentsSet(t, ints1, s1)
	assert.Equal(t, false, ok)

	// 0 is the value on the root node of s, which is kind of a special case. We
	// want to test deleting it and setting a new value (which should get put on
	// the root node).
	s2, ok := s.DelVal(0)
	assertSeqContentsSet(t, ints, s)
	assertSeqContentsSet(t, ints2, s2)
	assert.Equal(t, true, ok)

	s2, ok = s2.DelVal(0)
	assertSeqContentsSet(t, ints2, s2)
	assert.Equal(t, false, ok)

	s3, ok := s2.SetVal(5)
	assertSeqContentsSet(t, ints2, s2)
	assertSeqContentsSet(t, ints3, s3)
	assert.Equal(t, true, ok)
}

// Test getting values from a Set
func GetVal(t *T) {
	//Degenerate case
	s := NewSet()
	v, ok := s.GetVal(1)
	assert.Equal(t, nil, v)
	assert.Equal(t, false, ok)

	s = NewSet(0, 1, 2, 3, 4)
	v, ok = s.GetVal(1)
	assert.Equal(t, 1, t)
	assert.Equal(t, true, ok)

	// After delete
	s, _ = s.DelVal(1)
	v, ok = s.GetVal(1)
	assert.Equal(t, nil, v)
	assert.Equal(t, false, ok)

	// After set
	s, _ = s.SetVal(1)
	v, ok = s.GetVal(1)
	assert.Equal(t, 1, t)
	assert.Equal(t, true, ok)

	// After delete root node
	s, _ = s.DelVal(0)
	v, ok = s.GetVal(0)
	assert.Equal(t, nil, v)
	assert.Equal(t, false, ok)

	// After set root node
	s, _ = s.SetVal(5)
	v, ok = s.GetVal(5)
	assert.Equal(t, 5, t)
	assert.Equal(t, true, ok)
}

// Test that Size functions properly for all cases
func TestSetSize(t *T) {
	// Degenerate case
	s := NewSet()
	assert.Equal(t, uint64(0), s.Size())

	// Make sure setting on an empty Set produces the correct size
	s, _ = s.SetVal(0)
	assert.Equal(t, uint64(1), s.Size())

	// Initialization case
	s = NewSet(0, 1, 2)
	assert.Equal(t, uint64(3), s.Size())

	// Setting (both value not in and a value already in)
	s, _ = s.SetVal(3)
	assert.Equal(t, uint64(4), s.Size())
	s, _ = s.SetVal(3)
	assert.Equal(t, uint64(4), s.Size())

	// Deleting (both value already in and a value not in)
	s, _ = s.DelVal(3)
	assert.Equal(t, uint64(3), s.Size())
	s, _ = s.DelVal(3)
	assert.Equal(t, uint64(3), s.Size())

	// Deleting and setting the root node
	s, _ = s.DelVal(0)
	assert.Equal(t, uint64(2), s.Size())
	s, _ = s.SetVal(5)
	assert.Equal(t, uint64(3), s.Size())
}

// Test that Union functions properly
func TestUnion(t *T) {
	// Degenerate case
	empty := NewSet()
	assert.Equal(t, 0, Size(empty.Union(empty)))

	ints1 := []interface{}{0, 1, 2}
	ints2 := []interface{}{3, 4, 5}
	intsu := append(ints1, ints2...)
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)

	assertSeqContentsSet(t, ints1, s1.Union(empty))
	assertSeqContentsSet(t, ints1, empty.Union(s1))

	su := s1.Union(s2)
	assertSeqContentsSet(t, ints1, s1)
	assertSeqContentsSet(t, ints2, s2)
	assertSeqContentsSet(t, intsu, su)
}

// Test that Intersection functions properly
func TestIntersection(t *T) {
	// Degenerate case
	empty := NewSet()
	assert.Equal(t, 0, Size(empty.Intersection(empty)))

	ints1 := []interface{}{0, 1, 2}
	ints2 := []interface{}{1, 2, 3}
	ints3 := []interface{}{4, 5, 6}
	intsi := []interface{}{1, 2}
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)
	s3 := NewSet(ints3...)

	assert.Equal(t, 0, Size(s1.Intersection(empty)))
	assert.Equal(t, 0, Size(empty.Intersection(s1)))

	si := s1.Intersection(s2)
	assert.Equal(t, 0, Size(s1.Intersection(s3)))
	assertSeqContentsSet(t, ints1, s1)
	assertSeqContentsSet(t, ints2, s2)
	assertSeqContentsSet(t, ints3, s3)
	assertSeqContentsSet(t, intsi, si)
}

// Test that Difference functions properly
func TestDifference(t *T) {
	// Degenerate case
	empty := NewSet()
	assert.Equal(t, 0, Size(empty.Difference(empty)))

	ints1 := []interface{}{0, 1, 2, 3}
	ints2 := []interface{}{2, 3, 4}
	intsd := []interface{}{0, 1}
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)

	assertSeqContentsSet(t, ints1, s1.Difference(empty))
	assert.Equal(t, 0, Size(empty.Difference(s1)))

	sd := s1.Difference(s2)
	assertSeqContentsSet(t, ints1, s1)
	assertSeqContentsSet(t, ints2, s2)
	assertSeqContentsSet(t, intsd, sd)
}

// Test that SymDifference functions properly
func TestSymDifference(t *T) {
	// Degenerate case
	empty := NewSet()
	assert.Equal(t, 0, Size(empty.SymDifference(empty)))

	ints1 := []interface{}{0, 1, 2, 3}
	ints2 := []interface{}{2, 3, 4}
	intsd := []interface{}{0, 1, 4}
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)

	assertSeqContentsSet(t, ints1, s1.SymDifference(empty))
	assertSeqContentsSet(t, ints1, empty.SymDifference(s1))

	sd := s1.SymDifference(s2)
	assertSeqContentsSet(t, ints1, s1)
	assertSeqContentsSet(t, ints2, s2)
	assertSeqContentsSet(t, intsd, sd)
}
