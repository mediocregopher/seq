package seq

import (
	. "testing"
)

func kvints(kvs ...*KV) ([]*KV, []interface{}) {
	ints := make([]interface{}, len(kvs))
	for i := range kvs {
		ints[i] = kvs[i]
	}
	return kvs, ints
}

// Test creating a Set and calling the Seq interface methods on it
func TestHashMapSeq(t *T) {
	kvs, ints := kvints(
		KeyVal(1, "one"),
		KeyVal(2, "two"),
	)

	// Testing creation and Seq interface methods
	m := NewHashMap(kvs...)
	ms := testSeqNoOrderGen(t, m, ints)

	// ms should be empty at this point
	assertEmpty(ms, t)
}

// Test getting values from a HashMap
func TestHashMapGet(t *T) {
	kvs := []*KV{
		KeyVal(1, "one"),
		KeyVal(2, "two"),
	}

	// Degenerate case
	m := NewHashMap()
	assertEmpty(m, t)
	v, ok := m.Get(1)
	assertValue(v, nil, t)
	assertValue(ok, false, t)

	m = NewHashMap(kvs...)
	v, ok = m.Get(1)
	assertSeqContentsHashMap(m, kvs, t)
	assertValue(v, "one", t)
	assertValue(ok, true, t)

	v, ok = m.Get(3)
	assertSeqContentsHashMap(m, kvs, t)
	assertValue(v, nil, t)
	assertValue(ok, false, t)
}

// Test setting values on a HashMap
func TestHashMapSet(t *T) {

	// Set on empty
	m := NewHashMap()
	m1, ok := m.Set(1, "one")
	assertEmpty(m, t)
	assertSeqContentsHashMap(m1, []*KV{KeyVal(1, "one")}, t)
	assertValue(ok, true, t)

	// Set on same key
	m2, ok := m1.Set(1, "wat")
	assertSeqContentsHashMap(m1, []*KV{KeyVal(1, "one")}, t)
	assertSeqContentsHashMap(m2, []*KV{KeyVal(1, "wat")}, t)
	assertValue(ok, false, t)

	// Set on second new key
	m3, ok := m2.Set(2, "two")
	assertSeqContentsHashMap(m2, []*KV{KeyVal(1, "wat")}, t)
	assertSeqContentsHashMap(m3, []*KV{KeyVal(1, "wat"), KeyVal(2, "two")}, t)
	assertValue(ok, true, t)

}

// Test deleting keys from sets
func TestHashMapDel(t *T) {

	kvs := []*KV{
		KeyVal(1, "one"),
		KeyVal(2, "two"),
		KeyVal(3, "three"),
	}
	kvs1 := []*KV{
		KeyVal(2, "two"),
		KeyVal(3, "three"),
	}

	// Degenerate case
	m := NewHashMap()
	m1, ok := m.Del(1)
	assertEmpty(m, t)
	assertEmpty(m1, t)
	assertValue(ok, false, t)

	// Delete actual key
	m = NewHashMap(kvs...)
	m1, ok = m.Del(1)
	assertSeqContentsHashMap(m, kvs, t)
	assertSeqContentsHashMap(m1, kvs1, t)
	assertValue(ok, true, t)

	// Delete it again!
	m2, ok := m1.Del(1)
	assertSeqContentsHashMap(m1, kvs1, t)
	assertSeqContentsHashMap(m2, kvs1, t)
	assertValue(ok, false, t)

}
