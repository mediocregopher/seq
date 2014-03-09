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
