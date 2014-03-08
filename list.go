package seq

// A List is an implementation of Seq in the form of a single-linked-list, and
// is used as the underlying structure for Seqs for most methods that return a
// Seq. It is probably the most efficient and simplest of the implementations.
// Even though, conceptually, all Seq operations return a new Seq, the old Seq
// can actually share nodes with the new Seq (if both are Lists), thereby saving
// memory and copies.
type List struct {
	el   interface{}
	next *List
}

// Returns a new List comprised of the given elements (or no elements, for an
// empty list)
func NewList(els ...interface{}) *List {
	elsl := len(els)
	if elsl == 0 {
		return nil
	}

	var cur *List
	for i := 0; i < elsl; i++ {
		cur = &List{els[elsl-i-1], cur}
	}
	return cur
}

// Implementation of FirstRest for Seq interface. Completes in O(1) time.
func (l *List) FirstRest() (interface{}, Seq, bool) {
	if l == nil {
		return nil, l, false
	} else {
		return l.el, l.next, true
	}
}

// Implementation of String for Stringer interface.
func (l *List) String() string {
	return ToString(l, "(", ")")
}

// Prepends the given element to the front of the list, returning a copy of the
// new list. Completes in O(1) time.
func (l *List) Prepend(el interface{}) *List {
	return &List{el, l}
}

// Prepends the argument Seq to the beginning of the callee List, returning a
// copy of the new List. Completes in O(N) time, N being the length of the
// argument Seq
func (l *List) PrependSeq(s Seq) *List {
	var first, cur, prev *List
	var el interface{}
	var ok bool
	for {
		el, s, ok = s.FirstRest()
		if !ok {
			break
		}
		cur = &List{el, nil}
		if first == nil {
			first = cur
		}
		if prev != nil {
			prev.next = cur
		}
		prev = cur
	}

	// prev will be nil if s is empty
	if prev == nil {
		return l
	}

	prev.next = l
	return first
}

// Appends the given element to the end of the List, returning a copy of the new
// List. While most methods on List don't actually copy much data, this one
// copies the entire list. Completes in O(N) time.
func (l *List) Append(el interface{}) *List {
	var first, cur, prev *List
	for l != nil {
		cur = &List{l.el, nil}
		if first == nil {
			first = cur
		}
		if prev != nil {
			prev.next = cur
		}
		prev = cur
		l = l.next
	}
	final := &List{el, nil}
	if prev == nil {
		return final
	}
	prev.next = final
	return first
}

// Returns the nth index element (starting at 0), with bool being false if i is
// out of bounds. Completes in O(N) time.
func (l *List) Nth(n uint64) (interface{}, bool) {
	var el interface{}
	var ok bool
	s := Seq(l)
	for i := uint64(0);; i++ {
		el, s, ok = s.FirstRest()
		if !ok {
			return nil, false
		} else if i == n {
			return el, true
		}
	}
}
