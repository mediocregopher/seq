package seq

import (
	"fmt"
	"bytes"
)

// The general interface which most operations will actually operate on. Acts as
// an interface onto any data structure
type Seq interface {

	// Returns the number of elements contained in the data structure. This call
	// will always be O(1), no matter the data structure.
	Size() uint64

	// Returns the "first" element in the data structure as well as a Seq
	// containing a copy of the rest of the elements in the data structure. The
	// "first" element can be random for structures which don't have a concept
	// of order (like Set). Calling FirstRest on an empty Seq (Size() == 0) will
	// return "first" as nil, the same empty Seq (TODO copy? I don't think it
	// matters), and false. The third return value is true in all other cases.
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
