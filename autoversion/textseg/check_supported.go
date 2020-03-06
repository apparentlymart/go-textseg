package textseg

import (
	"unicode"
)

const unsupportedVersionMessage = "text segmentation is not supported for current Go runtime Unicode version " + unicode.Version

var supportedUnicodeVersions = map[string]struct{}{
	"9.0.0":  struct{}{},
	"10.0.0": struct{}{},
	"11.0.0": struct{}{},
	"12.0.0": struct{}{},
}

func assertSupportedVersion() {
	if !Supported() {
		panic(unsupportedVersionMessage)
	}
}

// Supported is a function specific to the autoversion wrappers of textseg
// which returns true if the current Go runtime uses a version of Unicode
// that this package can support.
//
// If this function returns false, calls to other functions in this package
// may panic.
func Supported() bool {
	_, ok := supportedUnicodeVersions[unicode.Version]
	return ok
}
