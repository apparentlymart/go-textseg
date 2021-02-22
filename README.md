# Unicode Text Segmentation for Go

This is an implementation of the Unicode Text Segmentation specification for Go.
Specifically, it currently includes only the "grapheme cluster" segmentation
algorithm.

## Unicode Version Support

Each major version of Unicode includes a set of tables that define how each
codepoint participates in the segmentation algorithms. Therefore any caller
of this library must select a specific version of Unicode to support.

To allow for each caller to make that decision separately even though
multiple callers may coexist in the same program, there is a separate
major release of this module for each supported major Unicode version.
Therefore you can select the specific version you want by module
path. For example, to use the algorithm and tables defined by Unicode
version 13:

```
go get github.com/apparentlymart/go-textseg/v13
```

```go
import (
    "github.com/apparentlymart/go-textseg/v13/textseg"
)
```

However, each release of Go also includes some Unicode-version-specific
functionality and you may prefer to use the text segmentation definitions
that are relevant to the version of Unicode that your Go runtime is
using elsewhere. To enable that, this repository has a special separate
module which uses the current Go runtime version to select a suitable
versioned implementation automatically:

```
go get github.com/apparentlymart/go-textseg/autoversion
```

```go
import (
    "github.com/apparentlymart/go-textseg/autoversion/textseg"
)
```

**IMPORTANT:** This "autoversion" wrapper uses Go build tags to select
a `go-textseg` major version based on the current Go version. We use
this strategy to ensure that only one version of `go-textseg` will
be compiled into your program, but the downside is that `go-textseg`
must be updated for each new Go release. If you use this library in
your program, you will need to fetch a new version of it each time
you switch to a new version of Go, even if that version of Go does
not introduce a new Unicode version.

## Usage

The most important function in each `textseg` package is
`ScanGraphemeClusters`, which is a function compatible with the
signature of `bufio.Scanner` in the Go standard library. Each
time the `Scan` function is called, the function will produce one
full grapheme cluster.
