package main

// A possibly useful example of implementing Seq on top of an io.Reader

import (
	"strings"
	"strconv"
	"fmt"
	"github.com/mediocregopher/seq"
	"bufio"
	"io"
	"os"
)

// SeqIoReader is a simple wrapper around a bufio.Reader
type SeqIoReader struct {
	reader *bufio.Reader
	delim  byte
}

func NewSeqIoReader(read io.Reader, delim byte) *SeqIoReader {
	return &SeqIoReader{
		reader: bufio.NewReader(read),
		delim:  delim,
	}
}

// FirstRest is all that is needed to implement the Seq interface. Check out the
// docs on FirstRest for more info about the three return values
func (sir *SeqIoReader) FirstRest() (interface{}, seq.Seq, bool) {
	data, err := sir.reader.ReadString(sir.delim)
	if err != nil {
		return nil, nil, false
	}
	return data, sir, true
}

// Given a string, converts it to an integer, increments it, and returns the
// incremented number
func IncStr(i interface{}) interface{} {
	trimmed := strings.TrimRight(i.(string), "\n")
	ii, err := strconv.Atoi(trimmed)
	if err != nil {
		panic(err)
	}
	return ii + 1
}

func main() {
	// We open a file which presumably is full of newline delimited numbers, and
	// create a SeqIoReader out of it
	file, err := os.Open("/tmp/numbers")
	if err != nil {
		panic(err)
	}
	sir := NewSeqIoReader(file, '\n')

	// We can call any of the seq functions on our SeqIoReader. We could have
	// also called ToList or ToSet if we wanted to keep the full output for
	// later. As it is right now when the file is fully read from it can't
	// reread the data.
	incdnums := seq.Map(IncStr, sir)
	fmt.Println(incdnums)
}
