package seq

import (
	"fmt"
	"reflect"
	"hash/crc32"
)

// This is an implementation of a persistent tree, which will then be used as
// the basis for vectors, hash maps, and hash sets.

type Setable interface {

	// Returns an integer for the value. For two equivalent values (as defined
	// by ==) Hash(i) should always return the same number. For multiple values
	// of i, Hash should return different values if possible.
	Hash(uint32) uint32

	// Given an arbitrary value found in a Set, returns whether or not the two
	// are equal
	Equal(interface{}) bool
}

// Returns an arbitrary integer for the given value/iteration tuple
func hash(v interface{}, i uint32) uint32 {
	switch vt := v.(type) {

	case Setable:
		return vt.Hash(i) % ARITY

	case uint:
		return uint32(vt) % ARITY
	case uint8:
		return uint32(vt) % ARITY
	case uint32:
		return uint32(vt) % ARITY
	case uint64:
		return uint32(vt) % ARITY
	case int:
		return uint32(vt) % ARITY
	case int8:
		return uint32(vt) % ARITY
	case int16:
		return uint32(vt) % ARITY
	case int32:
		return uint32(vt) % ARITY
	case int64:
		return uint32(vt) % ARITY
	case float32:
		return uint32(vt) % ARITY
	case float64:
		return uint32(vt) % ARITY
	
	case string:
		return crc32.ChecksumIEEE([]byte(vt)) % ARITY

	case []byte:
		return crc32.ChecksumIEEE(vt) % ARITY

	default:
		err := fmt.Sprintf("%s not hashable", reflect.TypeOf(v))
		panic(err)
	}
}

// Returns whether two values (potentially Setable's) are equivalent
func equal(v1, v2 interface{}) bool {
	if v1t, ok := v1.(Setable); ok {
		return v1t.Equal(v2)
	} else if v2t, ok := v2.(Setable); ok {
		return v2t.Equal(v1)
	} else {
		return v1 == v2
	}
}

// The number of children each node in Set (implemented as a hash tree) can have
const ARITY = 32;

type Set struct {
	// The value being held
	val  interface{}

	// Whether or not the held value has been set yet. Needed because the value
	// could be nil
	full bool

	// Slice of kids of this node. Could be an empty slice
	kids []*Set
}

func NewSet(vals ...interface{}) *Set {
	if len(vals) == 0 {
		return nil
	}
	s := new(Set)
	for i := range vals {
		s.setValDirty(vals[i], 0)
	}
	return s
}

// Tries to set the val on this Set node, or initialize the kids slice if it
// can't. Returns whether or not the value was set and whether or not it was
// already set.
//
// dirty
func (s *Set) shallowTrySetOrInit(val interface{}) (bool, bool) {
	if !s.full {
		s.val = val
		s.full = true
		return true, false
	} else if equal(s.val, val) {
		s.val = val
		s.full = true
		return true, true
	} else if s.kids == nil {
		s.kids = make([]*Set, ARITY)
	}
	return false, false
}

// dirty (obviously)
func (s *Set) setValDirty(val interface{}, i uint32) {
	if ok, _ := s.shallowTrySetOrInit(val); ok {
		return
	}
	
	h := hash(val, i)
	if kid := s.kids[h]; kid != nil {
		kid.setValDirty(val, i + 1)
	} else {
		s.kids[h] = NewSet(val)
	}
}

func (s *Set) clone() *Set {
	var newkids []*Set
	if s.kids != nil {
		newkids = make([]*Set, ARITY)
		copy(newkids, s.kids)
	}
	cs := &Set{
		val:  s.val,
		full: s.full,
		kids: newkids,
	}
	return cs
}

func (s *Set) SetVal(val interface{}) (*Set, bool) {
	if s == nil {
		return NewSet(val), false
	}

	cs := s.clone()
	root := cs
	i := uint32(0)
	var found bool
	var newcs, kid *Set
	for {
		if ok, prev := cs.shallowTrySetOrInit(val); ok {
			found = prev
			break
		}

		h := hash(val, i)
		if kid = cs.kids[h]; kid != nil {
			newcs = kid.clone()
			cs.kids[h] = newcs
			cs = newcs
			i++
		} else {
			cs.kids[h] = NewSet(val)
			break
		}
	}
	return root, found
}

func (s *Set) internalDelVal(val interface{}, i uint32) (*Set, bool) {
	if s == nil {
		return nil, false
	} else if s.full && equal(val, s.val) {
		cs := s.clone()
		cs.val = nil
		cs.full = false
		return cs, true
	} else if s.kids == nil {
		return s, false
	}

	h := hash(val, i)
	if newkid, ok := s.kids[h].internalDelVal(val, i + 1); ok {
		cs := s.clone()
		cs.kids[h] = newkid
		return cs, true
	}
	return s, false
}

func (s *Set) DelVal(val interface{}) (*Set, bool) {
	return s.internalDelVal(val, 0)
}

func (s *Set) internalFirstRest() (interface{}, *Set, bool) {
	if s == nil {
		return nil, nil, false
	}

	if s.kids != nil {
		var el interface{}
		var rest *Set
		var ok bool
		for i := range s.kids {
			if el, rest, ok = s.kids[i].internalFirstRest(); ok {
				cs := s.clone()
				cs.kids[i] = rest
				return el, cs, true
			}
		}
	}

	// We're not nil, but we don't have a value and no kids had values. We might
	// as well be nil.
	if !s.full {
		return nil, nil, false
	}

	return s.val, nil, true
}

func (s *Set) FirstRest() (interface{}, Seq, bool) {
	el, restSet, ok := s.internalFirstRest()
	return el, Seq(restSet), ok
}

func (s *Set) String() string {
	return ToString(s, "#{", "}#")
}
