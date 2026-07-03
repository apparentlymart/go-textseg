package ucdparse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

type TestDataScanner struct {
	lineScan *bufio.Scanner
}

func NewTestDataScanner(r io.Reader) TestDataScanner {
	return TestDataScanner{
		lineScan: bufio.NewScanner(r),
	}
}

func (sc TestDataScanner) NextEntry() (TestDataEntry, error) {
	for {
		if !sc.lineScan.Scan() {
			if err := sc.lineScan.Err(); err != nil {
				return TestDataEntry{}, err
			}
			return TestDataEntry{}, io.EOF
		}
		raw := sc.lineScan.Bytes()
		raw, _, _ = bytes.Cut(raw, commentMarker)
		raw = bytes.TrimSpace(raw)
		if len(raw) == 0 {
			continue
		}
		wantSegs := bytes.Split(raw, splitMarker)

		entry := TestDataEntry{}
		for _, wantSeg := range wantSegs {
			charsHex := bytes.Split(wantSeg, nonsplitMarker)
			var currentSeg []byte
			for _, charHex := range charsHex {
				charHexStr := string(bytes.TrimSpace(charHex))
				if len(charHexStr) == 0 {
					continue
				}
				charInt, err := strconv.ParseInt(charHexStr, 16, 32)
				if err != nil {
					return TestDataEntry{}, fmt.Errorf("invalid character reference %q", charHexStr)
				}
				buf := make([]byte, 4)
				n := utf8.EncodeRune(buf, rune(charInt))
				buf = buf[:n]
				entry.Input = append(entry.Input, buf...)
				currentSeg = append(currentSeg, buf...)
			}
			if len(currentSeg) != 0 {
				entry.WantSegments = append(entry.WantSegments, currentSeg)
			}
		}
		return entry, nil
	}
}

type TestDataEntry struct {
	Input        []byte
	WantSegments [][]byte
}

var splitMarker = []byte("÷")
var nonsplitMarker = []byte("×")
