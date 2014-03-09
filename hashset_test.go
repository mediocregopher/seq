package seq

import (
	. "testing"
)

// Test creating a Seq and calling the Seq interface methods on it
func TestSetSeq(t *T) {
	ints := []interface{}{1, "a", 5.0}

	// Testing creation and Seq interface methods
	s := NewSet(ints...)
	ss := testSeqNoOrderGen(t, s, ints)

	// ss should be empty at this point
	s = ToSet(ss)
	var nilpointer *Set
	assertEmpty(s, t)
	assertValue(s, nilpointer, t)
	assertValue(len(ToSlice(s)), 0, t)
}

// Test setting a value on a Set
func TestSetVal(t *T) {
	ints := []interface{}{0, 1, 2, 3, 4}
	ints1 := []interface{}{0, 1, 2, 3, 4, 5}

	// Degenerate case
	s := NewSet()
	assertEmpty(s, t)
	s, ok := s.SetVal(0)
	assertSeqContentsNoOrder(s, []interface{}{0}, t)
	assertValue(ok, true, t)

	s = NewSet(ints...)
	s1, ok := s.SetVal(5)
	assertSeqContentsNoOrder(s, ints, t)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertValue(ok, true, t)

	s2, ok := s1.SetVal(5)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertSeqContentsNoOrder(s2, ints1, t)
	assertValue(ok, false, t)
}

// Test deleting a value from a Set
func TestDelVal(t *T) {
	ints := []interface{}{0, 1, 2, 3, 4}
	ints1 := []interface{}{0, 1, 2, 3}
	ints2 := []interface{}{1, 2, 3 ,4}
	ints3 := []interface{}{1, 2, 3, 4, 5}

	// Degenerate case
	s := NewSet()
	assertEmpty(s, t)
	v, s, ok := s.DelVal(0)
	assertEmpty(s, t)
	assertValue(ok, false, t)

	s = NewSet(ints...)
	v, s1, ok := s.DelVal(4)
	assertValue(v, 4, t)
	assertSeqContentsNoOrder(s, ints, t)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertValue(ok, true, t)

	v, s1, ok = s1.DelVal(4)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertValue(ok, false, t)

	// 0 is the value on the root node of s, which is kind of a special case. We
	// want to test deleting it and setting a new value (which should get put on
	// the root node).
	v, s2, ok := s.DelVal(0)
	assertValue(v, 0, t)
	assertSeqContentsNoOrder(s, ints, t)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertValue(ok, true, t)

	v, s2, ok = s2.DelVal(0)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertValue(ok, false, t)

	s3, ok := s2.SetVal(5)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertSeqContentsNoOrder(s3, ints3, t)
	assertValue(ok, true, t)
}
