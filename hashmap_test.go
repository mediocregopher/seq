package seq

import (
	. "testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 0, Size(ms))
}

// Test getting values from a HashMap
func TestHashMapGet(t *T) {
	kvs := []*KV{
		KeyVal(1, "one"),
		KeyVal(2, "two"),
	}

	// Degenerate case
	m := NewHashMap()
	assert.Equal(t, 0, Size(m))
	v, ok := m.Get(1)
	assert.Equal(t, nil, v)
	assert.Equal(t, false, ok)

	m = NewHashMap(kvs...)
	v, ok = m.Get(1)
	assertSeqContentsHashMap(t, kvs, m)
	assert.Equal(t, "one", v)
	assert.Equal(t, true, ok)

	v, ok = m.Get(3)
	assertSeqContentsHashMap(t, kvs, m)
	assert.Equal(t, nil, v)
	assert.Equal(t, false, ok)
}

// Test setting values on a HashMap
func TestHashMapSet(t *T) {

	// Set on empty
	m := NewHashMap()
	m1, ok := m.Set(1, "one")
	assert.Equal(t, 0, Size(m))
	assertSeqContentsHashMap(t, []*KV{KeyVal(1, "one")}, m1)
	assert.Equal(t, true, ok)

	// Set on same key
	m2, ok := m1.Set(1, "wat")
	assertSeqContentsHashMap(t, []*KV{KeyVal(1, "one")}, m1)
	assertSeqContentsHashMap(t, []*KV{KeyVal(1, "wat")}, m2)
	assert.Equal(t, false, ok)

	// Set on second new key
	m3, ok := m2.Set(2, "two")
	assertSeqContentsHashMap(t, []*KV{KeyVal(1, "wat")}, m2)
	assertSeqContentsHashMap(t, []*KV{KeyVal(1, "wat"), KeyVal(2, "two")}, m3)
	assert.Equal(t, true, ok)

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
	assert.Equal(t, 0, Size(m))
	assert.Equal(t, 0, Size(m1))
	assert.Equal(t, false, ok)

	// Delete actual key
	m = NewHashMap(kvs...)
	m1, ok = m.Del(1)
	assertSeqContentsHashMap(t, kvs, m)
	assertSeqContentsHashMap(t, kvs1, m1)
	assert.Equal(t, true, ok)

	// Delete it again!
	m2, ok := m1.Del(1)
	assertSeqContentsHashMap(t, kvs1, m1)
	assertSeqContentsHashMap(t, kvs1, m2)
	assert.Equal(t, false, ok)

}

// Test that two hashmaps compare equality correctly
func TestHashMapEqual(t *T) {
	// Degenerate case
	hm1, hm2 := NewHashMap(), NewHashMap()
	assert.Equal(t, true, hm1.Equal(hm2))
	assert.Equal(t, true, hm2.Equal(hm1))

	// False with different sizes
	hm1, _ = hm1.Set("one", 1)
	assert.Equal(t, false, hm1.Equal(hm2))
	assert.Equal(t, false, hm2.Equal(hm1))

	// False with same sizes
	hm2, _ = hm2.Set("two", 2)
	assert.Equal(t, false, hm1.Equal(hm2))
	assert.Equal(t, false, hm2.Equal(hm1))

	// Now true
	hm1, _ = hm1.Set("two", 2)
	hm2, _ = hm2.Set("one", 1)
	assert.Equal(t, true, hm1.Equal(hm2))
	assert.Equal(t, true, hm2.Equal(hm1))

	// False with embedded HashMap
	hm1, _ = hm1.Set(NewHashMap().Set("three", 3))
	assert.Equal(t, false, hm1.Equal(hm2))
	assert.Equal(t, false, hm2.Equal(hm1))

	// True with embedded set
	hm2, _ = hm2.Set(NewHashMap().Set("three", 3))
	assert.Equal(t, true, hm1.Equal(hm2))
	assert.Equal(t, true, hm2.Equal(hm1))

	// False with same key, different value
	hm1, _ = hm1.Set("four", 4)
	hm2, _ = hm2.Set("four", 5)
	assert.Equal(t, false, hm1.Equal(hm2))
	assert.Equal(t, false, hm2.Equal(hm1))
}
