package seq

type Lazy struct {
	this  interface{}
	next  *Lazy
	ok    bool
	ch    chan struct{}
}


// A Lazy is an implementation of a Seq which only actually evaluates its
// contents as those contents become needed. Lazys can be chained together, so
// if you have three steps in a pipeline there aren't two intermediate Seqs
// created, only the final resulting one.
func NewLazy(t Thunk) *Lazy {
	l := &Lazy{ch: make(chan struct{})}
	go func(){
		l.ch <- struct{}{}
		el, next, ok := t()
		l.this = el
		l.next = NewLazy(next)
		l.ok = ok
		close(l.ch)
	}()
	return l
}

// Implementation of FirstRest for Seq interface. Completes in O(1) time.
func (l *Lazy) FirstRest() (interface{}, Seq, bool) {
	if l == nil {
		return nil, l, false
	}

	// Reading from the channel tells the Lazy to populate the data and prepare
	// the next item in the seq, it closes the channel when it's done that.
	if _, ok := <-l.ch; ok {
		<-l.ch
	}

	if l.ok {
		return l.this, l.next, true
	} else {
		return nil, nil, false
	}
}

// Thunks are the building blocks a Lazy. A Thunk returns an element, another
// Thunk, and a boolean representing if the call yielded any results or if it
// was actually empty (true indicates it yielded results).
type Thunk func() (interface{}, Thunk, bool)

func filterThunk(fn func(interface{}) bool, s Seq) Thunk {
	return func() (interface{}, Thunk, bool) {
		for {
			el, ns, ok := s.FirstRest()
			if !ok {
				return nil, nil, false
			}

			if keep := fn(el); keep {
				return el, filterThunk(fn, ns), true
			} else {
				s = ns
			}
		}
	}
}

// Lazy implementation of Filter
func LFilter(fn func(interface{}) bool, s Seq) Seq {
	return NewLazy(filterThunk(fn, s))
}
