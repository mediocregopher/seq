package seq

import (
	. "testing"

	"github.com/stretchr/testify/assert"
)

// Tests the FirstRest, Size, and ToSlice methods of a Seq
func testSeqGen(t *T, s Seq, ints []interface{}) Seq {
	intsl := uint64(len(ints))
	for i := range ints {
		assertSaneList(ToList(s), t)
		assert.Equal(t, intsl-uint64(i), Size(s))
		assert.Equal(t, ints[i:], ToSlice(s))

		first, rest, ok := s.FirstRest()
		assert.Equal(t, true, ok)
		assert.Equal(t, ints[i], first)

		s = rest
	}
	return s
}

// Tests the FirstRest, Size, and ToSlice methods of an unordered Seq
func testSeqNoOrderGen(t *T, s Seq, ints []interface{}) Seq {
	intsl := uint64(len(ints))

	m := map[interface{}]bool{}
	for i := range ints {
		m[ints[i]] = true
	}

	for i := range ints {
		assertSaneList(ToList(s), t)
		assert.Equal(t, intsl-uint64(i), Size(s))
		assertSeqContentsNoOrderMap(t, m, s)

		first, rest, ok := s.FirstRest()
		assert.Equal(t, true, ok)
		_, ok = m[first]
		assert.Equal(t, true, ok)

		delete(m, first)
		s = rest
	}
	return s
}

// Test reversing a Seq
func TestReverse(t *T) {
	// Normal case
	intl := []interface{}{3, 2, 1}
	l := NewList(intl...)
	nl := Reverse(l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{1, 2, 3}, ToSlice(nl))

	// Degenerate case
	l = NewList()
	nl = Reverse(l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, 0, Size(nl))
}

func testMapGen(t *T, mapFn func(func(interface{}) interface{}, Seq) Seq) {
	fn := func(n interface{}) interface{} {
		return n.(int) + 1
	}

	// Normal case
	intl := []interface{}{1, 2, 3}
	l := NewList(intl...)
	nl := mapFn(fn, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{2, 3, 4}, ToSlice(nl))

	// Degenerate case
	l = NewList()
	nl = mapFn(fn, l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, 0, Size(nl))
}

// Test mapping over a Seq
func TestMap(t *T) {
	testMapGen(t, Map)
}

// Test lazily mapping over a Seq
func TestLMap(t *T) {
	testMapGen(t, LMap)
}

// Test reducing over a Seq
func TestReduce(t *T) {
	fn := func(acc, el interface{}) (interface{}, bool) {
		return acc.(int) + el.(int), false
	}

	// Normal case
	intl := []interface{}{1, 2, 3, 4}
	l := NewList(intl...)
	r := Reduce(fn, 0, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 10, r)

	// Short-circuit case
	fns := func(acc, el interface{}) (interface{}, bool) {
		return acc.(int) + el.(int), el.(int) > 2
	}
	r = Reduce(fns, 0, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 6, r)

	// Degenerate case
	l = NewList()
	r = Reduce(fn, 0, l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, 0, r)
}

// Test the Any function
func TestAny(t *T) {
	fn := func(el interface{}) bool {
		return el.(int) > 3
	}

	// Value found case
	intl := []interface{}{1, 2, 3, 4}
	l := NewList(intl...)
	r, ok := Any(fn, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 4, r)
	assert.Equal(t, true, ok)

	// Value not found case
	intl = []interface{}{1, 2, 3}
	l = NewList(intl...)
	r, ok = Any(fn, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, nil, r)
	assert.Equal(t, false, ok)

	// Degenerate case
	l = NewList()
	r, ok = Any(fn, l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, nil, r)
	assert.Equal(t, false, ok)
}

// Test the All function
func TestAll(t *T) {
	fn := func(el interface{}) bool {
		return el.(int) > 3
	}

	// All match case
	intl := []interface{}{4, 5, 6}
	l := NewList(intl...)
	ok := All(fn, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, true, ok)

	// Not all match case
	intl = []interface{}{3, 4, 2, 5}
	l = NewList(intl...)
	ok = All(fn, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, false, ok)

	// Degenerate case
	l = NewList()
	ok = All(fn, l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, true, ok)
}

func testFilterGen(t *T, filterFn func(func(interface{}) bool, Seq) Seq) {
	fn := func(el interface{}) bool {
		return el.(int)%2 != 0
	}

	// Normal case
	intl := []interface{}{1, 2, 3, 4, 5}
	l := NewList(intl...)
	r := filterFn(fn, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{1, 3, 5}, ToSlice(r))

	// Degenerate cases
	l = NewList()
	r = filterFn(fn, l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, 0, Size(r))
}

// Test the Filter function
func TestFilter(t *T) {
	testFilterGen(t, Filter)
}

// Test the lazy Filter function
func TestLFilter(t *T) {
	testFilterGen(t, LFilter)
}

// Test Flatten-ing of a Seq
func TestFlatten(t *T) {
	// Normal case
	intl1 := []interface{}{0, 1, 2}
	intl2 := []interface{}{3, 4, 5}
	l1 := NewList(intl1...)
	l2 := NewList(intl2...)
	blank := NewList()
	intl := []interface{}{-1, l1, l2, 6, blank, 7}
	l := NewList(intl...)
	nl := Flatten(l)
	assert.Equal(t, intl1, ToSlice(l1))
	assert.Equal(t, intl2, ToSlice(l2))
	assert.Equal(t, 0, Size(blank))
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{-1, 0, 1, 2, 3, 4, 5, 6, 7}, ToSlice(nl))

	// Degenerate case
	nl = Flatten(blank)
	assert.Equal(t, 0, Size(blank))
	assert.Equal(t, 0, Size(nl))
}

func testTakeGen(t *T, takeFn func(uint64, Seq) Seq) {
	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4}
	l := NewList(intl...)
	nl := takeFn(3, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{0, 1, 2}, ToSlice(nl))

	// Edge cases
	nl = takeFn(5, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, intl, ToSlice(nl))

	nl = takeFn(6, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, intl, ToSlice(nl))

	// Degenerate cases
	empty := NewList()
	nl = takeFn(1, empty)
	assert.Equal(t, 0, Size(empty))
	assert.Equal(t, 0, Size(nl))

	nl = takeFn(0, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 0, Size(nl))
}

// Test taking from a Seq
func TestTake(t *T) {
	testTakeGen(t, Take)
}

// Test lazily taking from a Seq
func TestLTake(t *T) {
	testTakeGen(t, LTake)
}

func testTakeWhileGen(t *T, takeWhileFn func(func(interface{}) bool, Seq) Seq) {
	pred := func(el interface{}) bool {
		return el.(int) < 3
	}

	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4, 5}
	l := NewList(intl...)
	nl := takeWhileFn(pred, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{0, 1, 2}, ToSlice(nl))

	// Edge cases
	intl = []interface{}{5, 5, 5}
	l = NewList(intl...)
	nl = takeWhileFn(pred, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 0, Size(nl))

	intl = []interface{}{0, 1, 2}
	l = NewList(intl...)
	nl = takeWhileFn(pred, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{0, 1, 2}, ToSlice(nl))

	// Degenerate case
	l = NewList()
	nl = takeWhileFn(pred, l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, 0, Size(nl))
}

// Test taking from a Seq until a given condition
func TestTakeWhile(t *T) {
	testTakeWhileGen(t, TakeWhile)
}

// Test lazily taking from a Seq until a given condition
func TestLTakeWhile(t *T) {
	testTakeWhileGen(t, LTakeWhile)
}

// Test dropping from a Seq
func TestDrop(t *T) {
	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4}
	l := NewList(intl...)
	nl := Drop(3, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{3, 4}, ToSlice(nl))

	// Edge cases
	nl = Drop(5, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 0, Size(nl))

	nl = Drop(6, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 0, Size(nl))

	// Degenerate cases
	empty := NewList()
	nl = Drop(1, empty)
	assert.Equal(t, 0, Size(empty))
	assert.Equal(t, 0, Size(nl))

	nl = Drop(0, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, intl, ToSlice(nl))
}

// Test dropping from a Seq until a given condition
func TestDropWhile(t *T) {
	pred := func(el interface{}) bool {
		return el.(int) < 3
	}

	// Normal case
	intl := []interface{}{0, 1, 2, 3, 4, 5}
	l := NewList(intl...)
	nl := DropWhile(pred, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, []interface{}{3, 4, 5}, ToSlice(nl))

	// Edge cases
	intl = []interface{}{5, 5, 5}
	l = NewList(intl...)
	nl = DropWhile(pred, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, intl, ToSlice(nl))

	intl = []interface{}{0, 1, 2}
	l = NewList(intl...)
	nl = DropWhile(pred, l)
	assert.Equal(t, intl, ToSlice(l))
	assert.Equal(t, 0, Size(nl))

	// Degenerate case
	l = NewList()
	nl = DropWhile(pred, l)
	assert.Equal(t, 0, Size(l))
	assert.Equal(t, 0, Size(nl))
}
