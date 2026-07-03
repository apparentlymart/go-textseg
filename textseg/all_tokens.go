package textseg

import (
	"bufio"
	"bytes"
)

// AllTokens is a utility that uses a [bufio.SplitFunc] to produce a slice of
// all of the recognized tokens in the given buffer.
//
// Deprecated: This will be removed in a future release, because it's not
// strictly relevant to Unicode text segmentation.
func AllTokens(buf []byte, splitFunc bufio.SplitFunc) ([][]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(splitFunc)
	var ret [][]byte
	for scanner.Scan() {
		ret = append(ret, scanner.Bytes())
	}
	return ret, scanner.Err()
}

// TokenCount is a utility that uses a [bufio.SplitFunc] to count the number of
// recognized tokens in the given buffer.
//
// Deprecated: This will be removed in a future release, because it's not
// strictly relevant to Unicode text segmentation.
func TokenCount(buf []byte, splitFunc bufio.SplitFunc) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(splitFunc)
	var ret int
	for scanner.Scan() {
		ret++
	}
	return ret, scanner.Err()
}
