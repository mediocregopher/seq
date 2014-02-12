package seq

import (
	"bytes"
	"fmt"
	"testing"
)

// Returns whether or not two interface{} slices contain the same elements
func intSlicesEq(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Asserts that the given Seq is empty (contains no elements)
func assertEmpty(s Seq, t *testing.T) {
	if s.Size() != 0 {
		t.Fatalf("Seq isn't empty: %v", s.ToSlice())
	}
}

// Asserts that the given Seq has the given elements
func assertSeqContents(s Seq, intl []interface{}, t *testing.T) {
	if ls := s.ToSlice(); !intSlicesEq(ls, intl) {
		t.Fatalf("Slice contents wrong: %v not %v", ls, intl)
	}
}

// Asserts that v1 is the same as v2
func assertValue(v1, v2 interface{}, t *testing.T) {
	if v1 != v2 {
		t.Fatalf("Value wrong: %v not %v", v1, v2)
	}
}

// Turns a Seq into a string, with each element separated by a space and with a
// dstart and dend wrapping the whole thing
func toString(s Seq, dstart, dend string) string {
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
