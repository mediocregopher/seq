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
	s, ok := s.DelVal(0)
	assertEmpty(s, t)
	assertValue(ok, false, t)

	s = NewSet(ints...)
	s1, ok := s.DelVal(4)
	assertSeqContentsNoOrder(s, ints, t)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertValue(ok, true, t)

	s1, ok = s1.DelVal(4)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertValue(ok, false, t)

	// 0 is the value on the root node of s, which is kind of a special case. We
	// want to test deleting it and setting a new value (which should get put on
	// the root node).
	s2, ok := s.DelVal(0)
	assertSeqContentsNoOrder(s, ints, t)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertValue(ok, true, t)

	s2, ok = s2.DelVal(0)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertValue(ok, false, t)

	s3, ok := s2.SetVal(5)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertSeqContentsNoOrder(s3, ints3, t)
	assertValue(ok, true, t)
}

// Test getting values from a Set
func GetVal(t *T) {
	//Degenerate case
	s := NewSet()
	v, ok := s.GetVal(1)
	assertValue(v, nil, t)
	assertValue(ok, false, t)

	s = NewSet(0, 1, 2, 3, 4)
	v, ok = s.GetVal(1)
	assertValue(v, 1, t)
	assertValue(ok, true, t)

	// After delete
	s, _ = s.DelVal(1)
	v, ok = s.GetVal(1)
	assertValue(v, nil, t)
	assertValue(ok, false, t)

	// After set
	s, _ = s.SetVal(1)
	v, ok = s.GetVal(1)
	assertValue(v, 1, t)
	assertValue(ok, true, t)

	// After delete root node
	s, _ = s.DelVal(0)
	v, ok = s.GetVal(0)
	assertValue(v, nil, t)
	assertValue(ok, false, t)

	// After set root node
	s, _ = s.SetVal(5)
	v, ok = s.GetVal(5)
	assertValue(v, 5, t)
	assertValue(ok, true, t)
}

// Test that Size functions properly for all cases
func TestSetSize(t *T) {
	// Degenerate case
	s := NewSet()
	assertValue(s.Size(), uint64(0), t)

	// Initialization case
	s = NewSet(0, 1, 2)
	assertValue(s.Size(), uint64(3), t)

	// Setting (both value not in and a value already in)
	s, _ = s.SetVal(3)
	assertValue(s.Size(), uint64(4), t)
	s, _ = s.SetVal(3)
	assertValue(s.Size(), uint64(4), t)

	// Deleting (both value already in and a value not in)
	s, _ = s.DelVal(3)
	assertValue(s.Size(), uint64(3), t)
	s, _ = s.DelVal(3)
	assertValue(s.Size(), uint64(3), t)

	// Deleting and setting the root node
	s, _ = s.DelVal(0)
	assertValue(s.Size(), uint64(2), t)
	s, _ = s.SetVal(5)
	assertValue(s.Size(), uint64(3), t)

}

// Test that Union functions properly
func TestUnion(t *T) {
	// Degenerate case
	empty := NewSet()
	assertEmpty(empty.Union(empty), t)

	ints1 := []interface{}{0, 1, 2}
	ints2 := []interface{}{3, 4, 5}
	intsu := append(ints1, ints2...)
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)

	assertSeqContentsNoOrder(s1.Union(empty), ints1, t)
	assertSeqContentsNoOrder(empty.Union(s1), ints1, t)

	su := s1.Union(s2)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertSeqContentsNoOrder(su, intsu, t)
}

// Test that Intersection functions properly
func TestIntersection(t *T) {
	// Degenerate case
	empty := NewSet()
	assertEmpty(empty.Intersection(empty), t)

	ints1 := []interface{}{0, 1, 2}
	ints2 := []interface{}{1, 2, 3}
	ints3 := []interface{}{4, 5, 6}
	intsi := []interface{}{1, 2}
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)
	s3 := NewSet(ints3...)

	assertEmpty(s1.Intersection(empty), t)
	assertEmpty(empty.Intersection(s1), t)

	si := s1.Intersection(s2)
	assertEmpty(s1.Intersection(s3), t)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertSeqContentsNoOrder(s3, ints3, t)
	assertSeqContentsNoOrder(si, intsi, t)
}

// Test that Difference functions properly
func TestDifference(t *T) {
	// Degenerate case
	empty := NewSet()
	assertEmpty(empty.Difference(empty), t)

	ints1 := []interface{}{0, 1, 2, 3}
	ints2 := []interface{}{2, 3, 4}
	intsd := []interface{}{0, 1}
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)

	assertSeqContentsNoOrder(s1.Difference(empty), ints1, t)
	assertEmpty(empty.Difference(s1), t)

	sd := s1.Difference(s2)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertSeqContentsNoOrder(sd, intsd, t)
}

// Test that SymDifference functions properly
func TestSymDifference(t *T) {
	// Degenerate case
	empty := NewSet()
	assertEmpty(empty.SymDifference(empty), t)

	ints1 := []interface{}{0, 1, 2, 3}
	ints2 := []interface{}{2, 3, 4}
	intsd := []interface{}{0, 1, 4}
	s1 := NewSet(ints1...)
	s2 := NewSet(ints2...)

	assertSeqContentsNoOrder(s1.SymDifference(empty), ints1, t)
	assertSeqContentsNoOrder(empty.SymDifference(s1), ints1, t)

	sd := s1.SymDifference(s2)
	assertSeqContentsNoOrder(s1, ints1, t)
	assertSeqContentsNoOrder(s2, ints2, t)
	assertSeqContentsNoOrder(sd, intsd, t)
}
