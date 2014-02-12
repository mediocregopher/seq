package seq

// A List is an implementation of Seq in the form of a single-linked-list, and
// is used as the underlying structure for Seqs for most methods that return a
// Seq. It is probably the most efficient and simplest of the implementations.
// Even though, conceptually, all Seq operations return a new Seq, the old Seq
// can actually share nodes with the new Seq (if both are Lists), thereby saving
// memory and copies.
type List struct {
	el   interface{}
	size uint64
	next *List
}

// Returns a new List comprised of the given elements (or no elements, for an
// empty list)
func NewList(els ...interface{}) *List {
	elsl := uint64(len(els))
	zeroNode := &List{nil, 0, nil}
	if elsl == 0 {
		return zeroNode
	}

	var first, cur, prev *List
	for i := elsl; i > 0; i-- {
		cur = &List{els[elsl-i], i, nil}
		if first == nil {
			first = cur
		}
		if prev != nil {
			prev.next = cur
		}
		prev = cur
	}
	prev.next = zeroNode
	return first
}

// Implementation of Size for Seq interface. Completes in O(1) time.
func (l *List) Size() uint64 {
	return l.size
}

// Implementation of FirstRest for Seq interface. Completes in O(1) time.
func (l *List) FirstRest() (interface{}, Seq, bool) {
	el := l.el
	if l.size == 0 {
		return el, l, false
	} else {
		return el, l.next, true
	}
}

// Implementation of ToSlice for Seq interface. Completes in O(N) time.
func (l *List) ToSlice() []interface{} {
	ret := make([]interface{}, l.Size())
	s := Seq(l)
	i := 0
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			ret[i] = el
			i++
		} else {
			return ret
		}
	}
}

// Implementation of ToList for Seq interface. Completes in O(1) time.
func (l *List) ToList() *List {
	return l
}

// Implementation of String for Seq interface.
func (l *List) String() string {
	return toString(l, "(", ")")
}

// Prepends the given element to the front of the list, returning a copy of the
// new list. Completes in O(1) time.
func (l *List) Prepend(el interface{}) *List {
	return &List{el, l.size + 1, l}
}

// Prepends the argument Seq to the beginning of the callee List, returning a
// copy of the new List. Completes in O(N) time, N being the length of the
// argument Seq
func (l *List) PrependSeq(s Seq) *List {
	lsize := l.Size()
	ssize := s.Size()

	// Degenerate cases
	if lsize == 0 {
		return s.ToList()
	} else if ssize == 0 {
		return l
	}

	//Actual work, both the List and the Seq have elements
	var first, cur, prev *List
	var el interface{}
	var ok bool
	i := ssize
	for {
		if el, s, ok = s.FirstRest(); ok {
			cur = &List{el, i + lsize, nil}
			if first == nil {
				first = cur
			}
			if prev != nil {
				prev.next = cur
			}
			prev = cur
			i--
		} else {
			prev.next = l
			return first
		}
	}
}

// Appends the given element to the end of the List, returning a copy of the new
// List. Completes in O(N) time.
func (l *List) Append(el interface{}) *List {
	lsize := l.Size()

	//Degenerate case
	if lsize == 0 {
		return NewList(el)
	}

	//Actual work, the List has elements
	var first, cur, prev *List
	for {
		if l.next == nil {
			prev.next = NewList(el)
			return first
		}
		cur = &List{l.el, l.size + 1, nil}
		if first == nil {
			first = cur
		}
		if prev != nil {
			prev.next = cur
		}
		l = l.next
		prev = cur
	}
}

// Returns a reversed copy of the List. Completes in O(N) time.
func (l *List) Reverse() *List {
	nl := NewList()
	s := Seq(l)
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); ok {
			nl = nl.Prepend(el)
		} else {
			return nl
		}
	}
}

// Returns the nth index element (starting at 0), with bool being false if i is
// out of bounds. Completes in O(N) time.
func (l *List) Nth(i uint64) (interface{}, bool) {
	if i >= l.Size() {
		return nil, false
	}

	var el interface{}
	s := Seq(l)
	for j := uint64(0); ; j++ {
		el, s, _ = s.FirstRest()
		if j == i {
			return el, true
		}
	}
}
