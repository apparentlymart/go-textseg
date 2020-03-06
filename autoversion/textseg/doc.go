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
// The functions in this package will pacic if the Go runtime is implementing
// a Unicode version that does not yet have a textseg implementation, to
// avoid a situation where it might appear that versions are consistent when
// they actually are not.
// Therefore any caller of this package should expect to need to upgrade
// to a newer version of it each time the Go standard library switches to a
// newer Unicode version, and releases of textseg are likely to lag behind
// releases of Go itself.
package textseg
