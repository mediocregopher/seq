package main

// A possibly useful example of constructing our own lazy list

import (
	"fmt"
	"github.com/mediocregopher/seq"
	"io"
	"os"
	"bufio"
	"strings"
)

// Thunks are weird, but they are what's needed in order to create Lazys. This
// one is a wrapper around bufio.Reader
func seqIoThunk(reader *bufio.Reader, delim byte) seq.Thunk {
	return func() (interface{}, seq.Thunk, bool) {
		data, err := reader.ReadString(delim)
		if err != nil {
			return nil, nil, false
		}
		tdata := strings.TrimRight(data, "\n")
		return tdata, seqIoThunk(reader, delim), true
	}
}

func NewSeqIoReader(read io.Reader, delim byte) *seq.Lazy {
	thunk := seqIoThunk(bufio.NewReader(read), delim)
	return seq.NewLazy(thunk)
}

func main() {
	// We open a file which presumably is full of newline delimited numbers, and
	// create a SeqIoReader out of it
	file, err := os.Open("/tmp/numbers")
	if err != nil {
		panic(err)
	}
	sir := NewSeqIoReader(file, '\n')

	// The first time we call Println the file will be read fully. The second
	// time it's fully read, but lazy lists are cached so we can see the results
	// again
	fmt.Println(sir)
	fmt.Println(sir)
}
