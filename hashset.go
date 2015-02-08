package seq

import (
	"fmt"
	"hash/crc32"
	"reflect"
)

// This is an implementation of a persistent tree, which will then be used as
// the basis for vectors, hash maps, and hash sets.

// Setable is an interface for customizing how items being put in a Set are
// compared to each other to determine their equality
type Setable interface {

	// Returns an integer for the value. For two equivalent values (as defined
	// by Equal) Hash(i) should always return the same number. For multiple
	// values of i, Hash should return different values if possible.
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
	} else if v1t, ok := v1.([]byte); ok {
		if v2t, ok := v2.([]byte); ok {
			if len(v1t) != len(v2t) {
				return false
			}
			for i := range v1t {
				if v1t[i] != v2t[i] {
					return false
				}
			}
			return true
		}
		return false
	}
	return v1 == v2
}

// The number of children each node in Set (implemented as a hash tree) can have
const ARITY = 32

// A Set is an implementation of Seq in the form of a persistant hash-tree. All
// public operations on it return a new, immutable form of the modified
// variable, leaving the old one intact. Immutability is implemented through
// node sharing, so operations aren't actually copying the entire hash-tree
// everytime, only the nodes which change, making the implementation very
// efficient compared to just copying.
//
// Items in sets need to be hashable and comparable. This means they either need
// to be some real numeric type (int, float32, etc...), string, []byte, or
// implement the Setable interface.
type Set struct {

	// The value being held
	val interface{}

	// Whether or not the held value has been set yet. Needed because the value
	// could be nil
	full bool

	// Slice of kids of this node. Could be an empty slice
	kids []*Set

	// Number of values in this Set.
	size uint64
}

// NewSet returns a new Set of the given elements (or no elements, for an empty
// set)
func NewSet(vals ...interface{}) *Set {
	if len(vals) == 0 {
		return nil
	}
	set := new(Set)
	for i := range vals {
		set.setValDirty(vals[i], 0)
	}
	set.size = uint64(len(vals))
	return set
}

// Methods marked as "dirty" operate on the node in place, and potentially
// change it or its children.

// Dirty. Tries to set the val on this Set node, or initialize the kids slice if
// it can't. Returns whether or not the value was set and whether or not it was
// already set.
func (set *Set) shallowTrySetOrInit(val interface{}) (bool, bool) {
	if !set.full {
		set.val = val
		set.full = true
		return true, false
	} else if equal(set.val, val) {
		set.val = val
		set.full = true
		return true, true
	} else if set.kids == nil {
		set.kids = make([]*Set, ARITY)
	}
	return false, false
}

// dirty (obviously). Sets a value on this node in place. Only used during
// initialization.
func (set *Set) setValDirty(val interface{}, i uint32) {
	if ok, _ := set.shallowTrySetOrInit(val); ok {
		return
	}

	h := hash(val, i)
	if kid := set.kids[h]; kid != nil {
		kid.setValDirty(val, i+1)
	} else {
		set.kids[h] = NewSet(val)
	}
}

// Returns a copy of this set node, including allocating and copying the kids
// slice.
func (set *Set) clone() *Set {
	var newkids []*Set
	if set.kids != nil {
		newkids = make([]*Set, ARITY)
		copy(newkids, set.kids)
	}
	cs := &Set{
		val:  set.val,
		full: set.full,
		kids: newkids,
		size: set.size,
	}
	return cs
}

// The actual implementation of SetVal, because we need to pass i down the stack
func (set *Set) internalSetVal(val interface{}, i uint32) (*Set, bool) {
	if set == nil {
		return NewSet(val), true
	}
	cset := set.clone()
	if ok, prev := cset.shallowTrySetOrInit(val); ok {
		return cset, !prev
	}

	h := hash(val, i)
	newkid, ok := cset.kids[h].internalSetVal(val, i+1)
	cset.kids[h] = newkid
	return cset, ok
}

// SetVal returns a new Set with the given value added to it. Also returns
// whether or not this is the first time setting this value (false if it was
// already there and was overwritten). Completes in O(log(N)) time.
func (set *Set) SetVal(val interface{}) (*Set, bool) {
	nset, ok := set.internalSetVal(val, 0)
	if ok {
		if set == nil {
			nset.size = 1
		} else {
			nset.size = set.size + 1
		}
	}
	return nset, ok
}

// The actual implementation of DelVal, because we need to pass i down the stack
func (set *Set) internalDelVal(val interface{}, i uint32) (*Set, bool) {
	if set == nil {
		return nil, false
	} else if set.full && equal(val, set.val) {
		cset := set.clone()
		cset.val = nil
		cset.full = false
		return cset, true
	} else if set.kids == nil {
		return set, false
	}

	h := hash(val, i)
	if newkid, ok := set.kids[h].internalDelVal(val, i+1); ok {
		cset := set.clone()
		cset.kids[h] = newkid
		return cset, true
	}
	return set, false
}

// DelVal returns a new Set with the given value removed from it and whether or
// not the value was actually removed. Completes in O(log(N)) time.
func (set *Set) DelVal(val interface{}) (*Set, bool) {
	nset, ok := set.internalDelVal(val, 0)
	if ok && nset != nil {
		nset.size--
	}
	return nset, ok
}

// The actual implementation of GetVal, because we need to pass i down the stack
func (set *Set) internalGetVal(val interface{}, i uint32) (interface{}, bool) {
	if set == nil {
		return nil, false
	} else if set.full && equal(val, set.val) {
		return set.val, true
	} else if set.kids == nil {
		return nil, false
	}

	h := hash(val, i)
	return set.kids[h].internalGetVal(val, i+1)
}

// GetVal returns a value from the Set, along with  a boolean indiciating
// whether or not the value was found. Completes in O(log(N)) time.
func (set *Set) GetVal(val interface{}) (interface{}, bool) {
	return set.internalGetVal(val, 0)
}

// Actual implementation of FirstRest. Because we need it to return a *Set
// instead of Seq for one case.
func (set *Set) internalFirstRest() (interface{}, *Set, bool) {
	if set == nil {
		return nil, nil, false
	}

	if set.kids != nil {
		var el interface{}
		var rest *Set
		var ok bool
		for i := range set.kids {
			if el, rest, ok = set.kids[i].internalFirstRest(); ok {
				cset := set.clone()
				cset.kids[i] = rest
				return el, cset, true
			}
		}
	}

	// We're not nil, but we don't have a value and no kids had values. We might
	// as well be nil.
	if !set.full {
		return nil, nil, false
	}

	return set.val, nil, true
}

// FirstRest is an implementation of FirstRest for Seq interface. Completes in
// O(log(N)) time.
func (set *Set) FirstRest() (interface{}, Seq, bool) {
	el, restSet, ok := set.internalFirstRest()
	if ok && restSet != nil {
		restSet.size--
	}
	return el, Seq(restSet), ok
}

// Implementation of String for Stringer interface
func (set *Set) String() string {
	return ToString(set, "#{", "}#")
}

// Size returns the number of elements in the Set. Completes in O(1) time.
func (set *Set) Size() uint64 {
	if set == nil {
		return 0
	}
	return set.size
}

// Union returns a Set with all of the elements of the original Set along with
// everything in the given Seq. If an element is present in both the Set and the
// Seq, the element in the Seq overwrites. Completes in O(M*log(N)), with M
// being the number of elements in the Seq and N the number of elements in the
// Set
func (set *Set) Union(s Seq) *Set {
	if set == nil {
		return ToSet(s)
	}

	cset := set.clone()
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); !ok {
			return cset
		} else if cset, ok = cset.SetVal(el); ok {
			cset.size++
		}
	}
}

// Intersection returns a Set with all of the elements in Seq that are also in
// Set. Completes in O(M*log(N)), with M being the number of elements in the Seq
// and N the number of elements in the Set
func (set *Set) Intersection(s Seq) *Set {
	if set == nil {
		return nil
	}

	iset := NewSet()
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); !ok {
			return iset
		} else if _, ok = set.GetVal(el); ok {
			iset, _ = iset.SetVal(el)
		}
	}
}

// Difference returns a Set of all elements in the original Set that aren't in
// the Seq. Completes in O(M*log(N)), with M being the number of elements in the
// Seq and N the number of elements in the Set
func (set *Set) Difference(s Seq) *Set {
	if set == nil {
		return nil
	}

	cset := set.clone()
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); !ok {
			return cset
		}
		cset, _ = cset.DelVal(el)
	}
}

// SymDifference returns a Set of all elements that are either in the original
// Set or the given Seq, but not in both. Completes in O(M*log(N)), with M being
// the number of elements in the Seq and N the number of elements in the Set.
func (set *Set) SymDifference(s Seq) *Set {
	if set == nil {
		return ToSet(s)
	}

	cset := set.clone()
	var cset2 *Set
	var el interface{}
	var ok bool
	for {
		if el, s, ok = s.FirstRest(); !ok {
			return cset
		} else if cset2, ok = cset.DelVal(el); ok {
			cset = cset2
		} else {
			cset, _ = cset.SetVal(el)
		}
	}
}

// ToSet returns the elements in the Seq as a set. In general this completes in
// O(N*log(N)) time (I think...). If the given Seq is already a Set it will
// complete in O(1) time. If it is a HashMap it will complete in O(1) time, and
// the resultant Set will be comprised of all KVs
func ToSet(s Seq) *Set {
	if set, ok := s.(*Set); ok {
		return set
	} else if hm, ok := s.(*HashMap); ok {
		return hm.set
	}
	vals := ToSlice(s)
	return NewSet(vals...)
}
