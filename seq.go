package seq

import (
	"bytes"
	"fmt"
)

// Seq is the general interface which most operations will actually operate on.
// Acts as an interface onto any data structure
type Seq interface {

	// FirstRest returns the "first" element in the data structure as well as a
	// Seq containing a copy of the rest of the elements in the data structure.
	// The "first" element can be random for structures which don't have a
	// concept of order (like Set). Calling FirstRest on an empty Seq (Size() ==
	// 0) will return "first" as nil, the same empty Seq , and false. The third
	// return value is true in all other cases.
	FirstRest() (interface{}, Seq, bool)
}

// Size returns the number of elements contained in the data structure. In
// general this completes in O(N) time, except for Set and HashMap for which it
// completes in O(1)
func Size(s Seq) uint64 {
	switch st := s.(type) {
	case *Set:
		return st.Size()
	case *HashMap:
		return st.Size()
	default:
	}

	var ok bool
	for i := uint64(0); ; {
		if _, s, ok = s.FirstRest(); ok {
			i++
		} else {
			return i
		}
	}
}

// ToSlice returns the elements in the Seq as a slice. If the underlying Seq has
// any implicit order to it that order will be kept. An empty Seq will return an
// empty slice; nil is never returned. In general this completes in O(N) time.
func ToSlice(s Seq) []interface{} {
	var el interface{}
	var ok bool
	for ret := make([]interface{}, 0, 8); ; {
		if el, s, ok = s.FirstRest(); ok {
			ret = append(ret, el)
		} else {
			return ret
		}
	}
}

// ToString turns a Seq into a string, with each element separated by a space
// and with a dstart and dend wrapping the whole thing
func ToString(s Seq, dstart, dend string) string {
	buf := bytes.NewBufferString(dstart)
	buf.WriteString(" ")
	var el interface{}
	var strel fmt.Stringer
	var rest Seq
	var ok bool
	for {
		if el, rest, ok = s.FirstRest(); ok {
			if strel, ok = el.(fmt.Stringer); ok {
				buf.WriteString(strel.String())
			} else {
				buf.WriteString(fmt.Sprintf("%v", el))
			}
			buf.WriteString(" ")
			s = rest
		} else {
			break
		}
	}
	buf.WriteString(dend)
	return buf.String()
}

// Reverse returns a reversed copy of the List. Completes in O(N) time.
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

// Map returns a Seq consisting of the result of applying fn to each element in the
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

// ReduceFn is a function used in a reduce. The first argument is the
// accumulator, the second is an element from the Seq being reduced over. The
// ReduceFn returns the accumulator to be used in the next iteration, wherein
// that new accumulator will be called alongside the next element in the Seq.
// ReduceFn also returns a boolean representing whether or not the reduction
// should stop at this step. If true, the reductions will stop and any remaining
// elements in the Seq will be ignored.
type ReduceFn func(acc, el interface{}) (interface{}, bool)

// Reduce reduces over the given Seq using ReduceFn, with acc as the first
// accumulator value in the reduce. See ReduceFn for more details on how it
// works. The return value is the result of the reduction. Completes in O(N)
// time.
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

// Any returns the first element in Seq for which fn returns true, or nil. The
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

// All returns true if fn returns true for all elements in the Seq. Completes in
// O(N) time.
func All(fn func(interface{}) bool, s Seq) bool {
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			if !fn(el) {
				return false
			}
		} else {
			return true
		}
	}
}

// Filter returns a Seq containing all elements in the given Seq for which fn
// returned true. Completes in O(N) time.
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

// Flatten flattens the given Seq into a single, one-dimensional Seq. This
// method only flattens Seqs found in the top level of the given Seq, it does
// not recurse down to multiple layers. Completes in O(N*M) time, where N is the
// number of elements in the Seq and M is how large the Seqs in those elements
// actually are.
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

// Take returns a Seq containing the first n elements in the given Seq. If n is
// greater than the length of the given Seq then the whole Seq is returned.
// Completes in O(N) time.
func Take(n uint64, s Seq) Seq {
	l := NewList()
	var el interface{}
	var ok bool
	for i := uint64(0); i < n; i++ {
		el, s, ok = s.FirstRest()
		if !ok {
			break
		}
		l = l.Prepend(el)
	}
	return Reverse(l)
}

// TakeWhile goes through each item in the given Seq until an element returns
// false from pred. Returns a new Seq containing these truthful elements.
// Completes in O(N) time.
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

// Drop returns a Seq which the is the previous Seq without the first n
// elements. If n is greater than the length of the Seq, returns an empty Seq.
// Completes in O(N) time.
func Drop(n uint64, s Seq) Seq {
	var ok bool
	for i := uint64(0); i < n; i++ {
		_, s, ok = s.FirstRest()
		if !ok {
			break
		}
	}
	return s
}

// DropWhile drops elements from the given Seq until pred returns false for an
// element.  Returns a Seq of the remaining elements (including the one which
// returned false). Completes in O(N) time.
func DropWhile(pred func(interface{}) bool, s Seq) Seq {
	var el interface{}
	var curs Seq
	var ok bool
	for {
		el, curs, ok = s.FirstRest()
		if !ok || !pred(el) {
			break
		}
		s = curs
	}
	return s
}
