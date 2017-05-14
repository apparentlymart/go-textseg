package textseg

import (
	"bufio"
	"bytes"
)

// AllTokens is a utility that uses a bufio.SplitFunc to produce a slice of
// all of the recognized tokens in the given buffer.
func AllTokens(buf []byte, splitFunc bufio.SplitFunc) ([][]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(splitFunc)
	var ret [][]byte
	for scanner.Scan() {
		ret = append(ret, scanner.Bytes())
	}
	return ret, scanner.Err()
}
