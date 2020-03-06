// +build go1.13,!go1.14

package textseg

import (
	v11 "github.com/apparentlymart/go-textseg/v11/textseg"
)

// ScanGraphemeClusters is a split function for bufio.Scanner that splits on
// grapheme cluster boundaries, using the text segmentation rules from the
// Unicode version selected by the current Go runtime library.
//
// This function will appear to be missing if your current Go version is not
// supported by your current version of this package.
func ScanGraphemeClusters(data []byte, atEOF bool) (int, []byte, error) {
	return v11.ScanGraphemeClusters(data, atEOF)
}

// UnicodeMajorVersion is the major version of Unicode being used by this
// package. This should always match the first portion of the string returned
// by unicode.Version in the Go standard library.
const UnicodeMajorVersion = 11
