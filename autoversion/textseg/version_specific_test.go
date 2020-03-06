package textseg

// Note that due to the nature of this package this will only test with the
// go-textseg version that matches with the current Go runtime.

import (
	"testing"
)

func TestScanGraphemeClusters(t *testing.T) {
	// Our goal here is only to test that we're really calling into an
	// upstream ScanGraphemeClusters function and not, say, a different scanner
	// function by mistake. It's not intended to be a deep test and should
	// hopefully remain valid in future Unicode versions.
	t.Logf("testing with implementation for Unicode major version %d", UnicodeMajorVersion)

	const input = `whelp 🤦🏽‍♂️`
	got, err := TokenCount([]byte(input), ScanGraphemeClusters)
	want := 7
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("wrong number of tokens\ninput: %s\ngot:   %d\nwant:  %d", input, got, want)
	}
}
