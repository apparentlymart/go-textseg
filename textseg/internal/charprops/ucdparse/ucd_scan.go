package ucdparse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"iter"
	"strconv"
)

type Scanner struct {
	lineScan *bufio.Scanner
}

func NewScanner(r io.Reader) Scanner {
	return Scanner{
		lineScan: bufio.NewScanner(r),
	}
}

func (sc Scanner) NextEntry() (Entry, error) {
	for {
		if !sc.lineScan.Scan() {
			if err := sc.lineScan.Err(); err != nil {
				return Entry{}, err
			}
			return Entry{}, io.EOF
		}
		raw := sc.lineScan.Bytes()
		raw, _, _ = bytes.Cut(raw, commentMarker)
		rng, fieldsRaw, ok := bytes.Cut(raw, fieldSeparator)
		if !ok {
			continue // Not a data line, then. We'll ignore it.
		}
		rng = bytes.TrimSpace(rng)
		fieldsRaw = bytes.TrimSpace(fieldsRaw)
		startHex, endHex, isRange := bytes.Cut(rng, rangeSeparator)
		var start, end rune
		startInt, err := strconv.ParseUint(string(startHex), 16, 32)
		if err != nil {
			return Entry{}, fmt.Errorf("invalid start of character range %q: %s", startHex, err)
		}
		start = rune(startInt)
		if isRange {
			endInt, err := strconv.ParseUint(string(endHex), 16, 32)
			if err != nil {
				return Entry{}, fmt.Errorf("invalid end of character range %q: %s", endHex, err)
			}
			end = rune(endInt)
		} else {
			end = start
		}
		return Entry{
			Start:  start,
			End:    end + 1,
			Fields: Fields(fieldsRaw),
		}, nil
	}
}

type Entry struct {
	// Start is inclusive, while End is exclusive, as is typical for Go slicing.
	Start, End rune
	Fields     Fields
}

func (e Entry) FirstField() string {
	return e.Fields.FirstString()
}

func (e Entry) AllFields() iter.Seq[string] {
	return e.Fields.Strings()
}

type Fields []byte

func (f Fields) Strings() iter.Seq[string] {
	remain := []byte(f)
	return func(yield func(string) bool) {
		for {
			before, after, foundSep := bytes.Cut(remain, fieldSeparator)
			val := bytes.TrimSpace(before)
			if !yield(string(val)) {
				return
			}
			if !foundSep {
				return // all done, then
			}
			remain = after
		}
	}
}

func (f Fields) FirstString() string {
	before, _, _ := bytes.Cut(f, fieldSeparator)
	return string(bytes.TrimSpace(before))
}

var commentMarker = []byte{'#'}
var fieldSeparator = []byte{';'}
var rangeSeparator = []byte{'.', '.'}
