package textseg

import (
	"unicode"

	v10 "github.com/apparentlymart/go-textseg/v10/textseg"
	v11 "github.com/apparentlymart/go-textseg/v11/textseg"
	v12 "github.com/apparentlymart/go-textseg/v12/textseg"
	v9 "github.com/apparentlymart/go-textseg/v9/textseg"
)

// ScanGraphemeClusters is a split function for bufio.Scanner that splits on
// grapheme cluster boundaries, using the text segmentation rules from the
// Unicode version selected by the current Go runtime library.
//
// This function will panic if the Go runtime is using a version that this
// package does not support.
func ScanGraphemeClusters(data []byte, atEOF bool) (int, []byte, error) {
	switch unicode.Version {
	case "9.0.0":
		return v9.ScanGraphemeClusters(data, atEOF)
	case "10.0.0":
		return v10.ScanGraphemeClusters(data, atEOF)
	case "11.0.0":
		return v11.ScanGraphemeClusters(data, atEOF)
	case "12.0.0":
		return v12.ScanGraphemeClusters(data, atEOF)
	default:
		panic(unsupportedVersionMessage)
	}
}
