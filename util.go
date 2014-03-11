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

// Asserts that the given Seq has all elements, and only the elements, in the
// given map
func assertSeqContentsNoOrderMap(s Seq, m map[interface{}]bool, t *testing.T) {
	ls := ToSlice(s)
	if len(ls) != len(m) {
		t.Fatalf("Slice contents wrong: %v not %v", ls, m)
	}
	for i := range ls {
		if _, ok := m[ls[i]]; !ok {
			t.Fatalf("Slice contents wrong: %v not %v", ls, m)
		}
	}
}

// Asserts that the given Seq has all the elements, and only the elements
// (duplicates removed), in the given slice, although no necessarily in the
// order given in the slice
func assertSeqContentsSet(s Seq, ints []interface{}, t *testing.T) {
	m := map[interface{}]bool{}
	for i := range ints {
		m[ints[i]] = true
	}
	assertSeqContentsNoOrderMap(s, m, t)
}

func assertSeqContentsHashMap(s Seq, kvs []*KV, t *testing.T) {
	m := map[interface{}]bool{}
	for i := range kvs {
		m[*kvs[i]] = true
	}
	ls := ToSlice(s)
	if len(ls) != len(m) {
		t.Fatalf("Slice contents wrong: %v not %v", ls, m)
	}
	for i := range ls {
		kv := ls[i].(*KV)
		if _, ok := m[*kv]; !ok {
			t.Fatalf("Slice contents wrong: %v not %v", ls, m)
		}
	}
}

// Asserts that v1 is the same as v2
func assertValue(v1, v2 interface{}, t *testing.T) {
	if v1 != v2 {
		t.Fatalf("Value wrong: %v not %v", v1, v2)
	}
}

// Asserts that v1 is a key in the given map
func assertInMap(v1 interface{}, m map[interface{}]bool, t *testing.T) {
	if _, ok := m[v1]; !ok {
		t.Fatalf("Value not in set: %v not in %v", v1, m)
	}
}
