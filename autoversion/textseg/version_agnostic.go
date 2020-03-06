package textseg

// The functions in this file are not currently specific to a given version of
// go-textseg and so we just call into the version from the latest upstream
// major version in all cases.

import (
	"bufio"

	realImpl "github.com/apparentlymart/go-textseg/v12/textseg"
)

// ScanUTF8Sequences is a split function for bufio.Scanner that splits on UTF8 sequence boundaries.
//
// This is included largely for completeness, since this behavior is already built in to Go when ranging over a string.
func ScanUTF8Sequences(data []byte, atEOF bool) (int, []byte, error) {
	return realImpl.ScanUTF8Sequences(data, atEOF)
}

// TokenCount is a utility that uses a bufio.SplitFunc to count the number of recognized tokens in the given buffer.
func TokenCount(buf []byte, splitFunc bufio.SplitFunc) (int, error) {
	return realImpl.TokenCount(buf, splitFunc)
}

// AllTokens is a utility that uses a bufio.SplitFunc to produce a slice of all of the recognized tokens in the given buffer.
func AllTokens(buf []byte, splitFunc bufio.SplitFunc) ([][]byte, error) {
	return realImpl.AllTokens(buf, splitFunc)
}
