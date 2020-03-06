// Package textseg is a wrapper around each of the the Unicode-version-specific
// textseg implementations that selects an implementation automatically based
// on the Unicode version of the Go standard library that it's being built
// against.
//
// The main textseg modules are designed with one major version per Unicode
// major version so that different callers in the same program can potentially
// use the segmentation rules from different versions of Unicode. However, in
// programs that are using textseg in conjunction with other Unicode-related
// functionality in the Go standard library, it could be desirable to keep the
// textseg version aligned with the Go library's Unicode version so that the
// total behavior is consistent with a single Unicode version.
//
// The API of this package is a subset of the API of the underlying
// version-specific packages, intended as a drop-in replacement for callers
// who are using that subset.
//
// In order to avoid linking in the data for all unicode versions into a
// program that imports the autoversion package, this package selects a Unicode
// version based on the Go compiler version that the program is being built
// with. That means that this package will need to be upgraded for each new
// Go version that is released, even if it doesn't introduce a new version of
// Unicode. There will be a compile-time error indicating that the
// ScanGraphemeClusters function is unavailable if you are using a Go version
// that is not supported by your current textseg autoversion package.
package textseg
