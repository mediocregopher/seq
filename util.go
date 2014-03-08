package seq

import (
	"testing"
)

// Returns whether or not two interface{} slices contain the same elements
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

// Asserts that the given Seq is empty (contains no elements)
func assertEmpty(s Seq, t *testing.T) {
	if Size(s) != 0 {
		t.Fatalf("Seq isn't empty: %v", ToSlice(s))
	}
}

// Asserts that the given Seq has the given elements
func assertSeqContents(s Seq, intl []interface{}, t *testing.T) {
	if ls := ToSlice(s); !intSlicesEq(ls, intl) {
		t.Fatalf("Slice contents wrong: %v not %v", ls, intl)
	}
}

// Asserts that v1 is the same as v2
func assertValue(v1, v2 interface{}, t *testing.T) {
	if v1 != v2 {
		t.Fatalf("Value wrong: %v not %v", v1, v2)
	}
}
