package seq

// The general interface which most operations will actually operate on. Acts as
// an interface onto any data structure
type Seq interface {

	// Returns the number of elements contained in the data structure.
	Size() uint64

	// Returns the "first" element in the data structure as well as a Seq
	// containing a copy of the rest of the elements in the data structure. The
	// "first" element can be random for structures which don't have a concept
	// of order (like Set). Calling FirstRest on an empty Seq (Size() == 0) will
	// return "first" as nil, the same empty Seq , and false. The third return
	// value is true in all other cases.
	FirstRest() (interface{}, Seq, bool)

	// Returns the elements in the Seq as a slice. If the underlying Seq has any
	// implicit order to it that order will be kept. An empty Seq will return an
	// empty slice; nil is never returned. The slice's capacity will match
	// whatever is returned by Size on the same Seq.
	ToSlice() []interface{}

	// Returns the elements in the Seq as a List. Has similar properties as
	// ToSlice.
	ToList() *List

	// Returns a string representation of the element. This is useful for
	// debugging, but not much else
	String() string
}

// Returns a reversed copy of the List. Completes in O(N) time.
func Reverse(s Seq) Seq {
	l := NewList()
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			l = l.Prepend(el)
		} else {
			return l
		}
	}
}

// Returns a Seq consisting of the result of applying fn to each element in the
// given Seq. Completes in O(N) time.
func Map(fn func(interface{}) interface{}, s Seq) Seq {
	l := NewList()
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			l = l.Prepend(fn(el))
		} else {
			break
		}
	}
	return Reverse(l)
}

// A function used in a reduce. The first argument is the accumulator, the
// second is an element from the Seq being reduced over. The ReduceFn returns
// the accumulator to be used in the next iteration, wherein that new
// accumulator will be called alongside the next element in the Seq. ReduceFn
// also returns a boolean representing whether or not the reduction should stop
// at this step. If true, the reductions will stop and any remaining elements in
// the Seq will be ignored.
type ReduceFn func(acc, el interface{}) (interface{}, bool)

// Reduces over the given Seq using ReduceFn, with acc as the first accumulator
// value in the reduce. See ReduceFn for more details on how it works. The
// return value is the result of the reduction. Completes in O(N) time.
func Reduce(fn ReduceFn, acc interface{}, s Seq) interface{} {
	var el interface{}
	var ok, stop bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			acc, stop = fn(acc, el)
			if stop {
				break
			}
		} else {
			break
		}
	}
	return acc
}

// Returns the first element in Seq for which fn returns true, or nil. The
// returned boolean indicates whether or not a matching element was found.
// Completes in O(N) time.
func Any(fn func(el interface{}) bool, s Seq) (interface{}, bool) {
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			if fn(el) {
				return el, true
			}
		} else {
			return nil, false
		}
	}
}

// Returns a Seq containing all elements in the given Seq for which fn returned
// true. Completes in O(N) time.
func Filter(fn func(el interface{}) bool, s Seq) Seq {
	l := NewList()
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			if fn(el) {
				l = l.Prepend(el)
			}
		} else {
			return Reverse(l)
		}
	}
}

// Flattens the given Seq into a single, one-dimensional Seq. This method only
// flattens Seqs found in the top level of the given Seq, it does not recurse
// down to multiple layers. Completes in O(N*M) time, where N is the number of
// elements in the Seq and M is how large the Seqs in those elements actually
// are.
func Flatten(s Seq) Seq {
	l := NewList()
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			if els, ok := el.(Seq); ok {
				l = l.PrependSeq(Reverse(els))
			} else {
				l = l.Prepend(el)
			}
		} else {
			return Reverse(l)
		}
	}
}

// Returns a Seq containing the first n elements in the given Seq. If n is
// greater than the length of the given Seq then the whole Seq is returned.
// Completes in O(N) time.
func Take(n uint64, s Seq) Seq {
	l := NewList()
	var el interface{}
	var ok bool
	for i := uint64(0); i < n; i++{
		el, s, ok = s.FirstRest()
		if !ok {
			break
		}
		l = l.Prepend(el)
	}
	return Reverse(l)
}

// Goes through each item in the given Seq until an element returns false from
// pred. Returns a new Seq containing these truthful elements. Completes in O(N)
// time.
func TakeWhile(pred func(interface{}) bool, s Seq) Seq {
	l := NewList()
	var el interface{}
	var ok bool
	for {
		el, s, ok = s.FirstRest()
		if !ok || !pred(el) {
			break
		}
		l = l.Prepend(el)
	}
	return Reverse(l)
}
