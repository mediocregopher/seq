package seq

import (
	"runtime/debug"
	"testing"
)

// Asserts that the given Seq has all elements, and only the elements, in the
// given map
func assertSeqContentsNoOrderMap(t *testing.T, m map[interface{}]bool, s Seq) {
	ls := ToSlice(s)
	if len(ls) != len(m) {
		t.Fatalf("Slice contents wrong: %v not %v\n%s", ls, m, debug.Stack())
	}
	for i := range ls {
		if _, ok := m[ls[i]]; !ok {
			t.Fatalf("Slice contents wrong: %v not %v\n%s", ls, m, debug.Stack())
		}
	}
}

// Asserts that the given Seq has all the elements, and only the elements
// (duplicates removed), in the given slice, although no necessarily in the
// order given in the slice
func assertSeqContentsSet(t *testing.T, ints []interface{}, s Seq) {
	m := map[interface{}]bool{}
	for i := range ints {
		m[ints[i]] = true
	}
	assertSeqContentsNoOrderMap(t, m, s)
}

func assertSeqContentsHashMap(t *testing.T, kvs []*KV, s Seq) {
	m := map[interface{}]bool{}
	for i := range kvs {
		m[*kvs[i]] = true
	}
	ls := ToSlice(s)
	if len(ls) != len(m) {
		t.Fatalf("Slice contents wrong: %v not %v\n%s", ls, m, debug.Stack())
	}
	for i := range ls {
		kv := ls[i].(*KV)
		if _, ok := m[*kv]; !ok {
			t.Fatalf("Slice contents wrong: %v not %v\n%s", ls, m, debug.Stack())
		}
	}
}
