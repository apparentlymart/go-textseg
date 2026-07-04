package textseg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/apparentlymart/go-textseg/v17/textseg/internal/charprops"
	"github.com/apparentlymart/go-textseg/v17/textseg/internal/charprops/ucdparse"
)

const graphemeBreakTestDataFile = "internal/charprops/ucd/auxiliary/GraphemeBreakTest.txt"

func TestScanGraphemeClusters(t *testing.T) {
	testDataFile, err := os.Open(graphemeBreakTestDataFile)
	if err != nil {
		t.Fatal(err)
	}
	testDataScanner := ucdparse.NewTestDataScanner(testDataFile)

	for {
		test, err := testDataScanner.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		t.Run(fmt.Sprintf("%x", test.Input), func(t *testing.T) {
			got, err := allTokens(test.Input, ScanGraphemeClusters)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !reflect.DeepEqual(got, test.WantSegments) {
				// Also get the rune values resulting from decoding utf8,
				// since they are generally easier to look up to figure out
				// what's failing.
				runes := make([]string, 0, len(test.Input))
				seqs := make([][]byte, 0, len(test.Input))
				categories := make([]string, 0, len(test.Input))
				buf := test.Input
				for len(buf) > 0 {
					r, size := utf8.DecodeRune(buf)
					runes = append(runes, fmt.Sprintf("0x%04x", r))
					seqs = append(seqs, buf[:size])
					props, _ := charprops.LookupFirstChar(buf)
					categories = append(categories, props.String())
					buf = buf[size:]
				}

				t.Errorf(
					"wrong result\ninput: %s\nutf8s: %s\nrunes: %s\ncats:  %s\ngot:   %s\nwant:  %s",
					formatBytes(test.Input),
					formatByteRanges(seqs),
					strings.Join(runes, " "),
					strings.Join(categories, " "),
					formatByteRanges(got),
					formatByteRanges(test.WantSegments),
				)
			}
		})
	}
}

// TestScanGraphemeClusters_partial is a variant of TestScanGraphemeClusters
// that makes sure the same logic works when data arrives in smaller chunks,
// such as streaming over a socket.
func TestScanGraphemeClusters_partial(t *testing.T) {
	testDataFile, err := os.Open(graphemeBreakTestDataFile)
	if err != nil {
		t.Fatal(err)
	}
	testDataScanner := ucdparse.NewTestDataScanner(testDataFile)

	for {
		test, err := testDataScanner.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		t.Run(fmt.Sprintf("%x", test.Input), func(t *testing.T) {
			r := &dripReader{test.Input}
			sc := bufio.NewScanner(r)
			sc.Split(ScanGraphemeClusters)
			var got [][]byte
			for sc.Scan() {
				got = append(got, sc.Bytes())
			}
			err := sc.Err()

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !reflect.DeepEqual(got, test.WantSegments) {
				// Also get the rune values resulting from decoding utf8,
				// since they are generally easier to look up to figure out
				// what's failing.
				runes := make([]string, 0, len(test.Input))
				seqs := make([][]byte, 0, len(test.Input))
				categories := make([]string, 0, len(test.Input))
				buf := test.Input
				for len(buf) > 0 {
					r, size := utf8.DecodeRune(buf)
					runes = append(runes, fmt.Sprintf("0x%04x", r))
					seqs = append(seqs, buf[:size])
					props, _ := charprops.LookupFirstChar(buf)
					categories = append(categories, props.String())
					buf = buf[size:]
				}

				t.Errorf(
					"wrong result\ninput: %s\nutf8s: %s\nrunes: %s\ncats:  %s\ngot:   %s\nwant:  %s",
					formatBytes(test.Input),
					formatByteRanges(seqs),
					strings.Join(runes, " "),
					strings.Join(categories, " "),
					formatByteRanges(got),
					formatByteRanges(test.WantSegments),
				)
			}
		})
	}
}

func TestScanGraphemeClusters_incompleteCluster(t *testing.T) {
	input := []byte("a")
	advance, token, err := ScanGraphemeClusters(input, false)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 0 {
		t.Fatalf("unexpected advance of %d with token %q; want zero, because we're not at a split point and not EOF", advance, token)
	}
	if token != nil {
		t.Fatalf("unexpected token %q; want nil", token)
	}

	// If we ask at EOF though, we should get told that the one remaining
	// character is a grapheme cluster.
	advance, token, err = ScanGraphemeClusters(input, true)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 1 {
		t.Fatalf("unexpected advance of %d with token %q; want 1", advance, token)
	}
	if string(token) != "a" {
		t.Fatalf("unexpected token %q; want \"a\"", token)
	}

	// We should also detect an incomplete cluster that's not at the beginning
	// of the string.
	input = []byte("\U00011F02\U00011F02")
	advance, token, err = ScanGraphemeClusters(input, false)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 0 {
		t.Fatalf("unexpected advance of %d with token %q; want zero, because we're not at a split point and not EOF", advance, token)
	}
	if token != nil {
		t.Fatalf("unexpected token %q; want nil", token)
	}

	// If we ask at EOF though, we should get told that the one remaining
	// character is a grapheme cluster.
	advance, token, err = ScanGraphemeClusters(input, true)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 8 {
		t.Fatalf("unexpected advance of %d with token %q; want 8", advance, token)
	}
	if string(token) != "\U00011F02\U00011F02" {
		t.Fatalf("unexpected token %q; want \"a\"", token)
	}
}

func TestScanGraphemeClusters_incompleteUTF8(t *testing.T) {
	input := []byte("\xc0")
	advance, token, err := ScanGraphemeClusters(input, false)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 0 {
		t.Fatalf("unexpected advance of %d with token %q; want zero, because we're not at a split point and not EOF", advance, token)
	}
	if token != nil {
		t.Fatalf("unexpected token %q; want nil", token)
	}

	// If we ask at EOF though, we should get told to consume this one byte
	// as a grapheme cluster even though it's invalid.
	advance, token, err = ScanGraphemeClusters(input, true)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 1 {
		t.Fatalf("unexpected advance of %d with token %q; want 1", advance, token)
	}
	if string(token) != "\xc0" {
		t.Fatalf("unexpected token %q; want \"\xc0\"", token)
	}

	// We should also detect an incomplete UTF-8 sequence that is not at the
	// beginning of a grapheme cluster.
	input = []byte("\U00011F02\xf0")
	advance, token, err = ScanGraphemeClusters(input, false)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 0 {
		t.Fatalf("unexpected advance of %d with token %q; want zero, because we're not at a split point and not EOF", advance, token)
	}
	if token != nil {
		t.Fatalf("unexpected token %q; want nil", token)
	}

	// If we ask at EOF though, we should get told to consume the first of
	// the two scalar values so that the invalid byte can be handled separately
	// by a subsequent call.
	advance, token, err = ScanGraphemeClusters(input, true)
	if err != nil {
		t.Fatal(err)
	}
	if advance != 4 {
		t.Fatalf("unexpected advance of %d with token %q; want 4", advance, token)
	}
	if string(token) != "\U00011F02" {
		t.Fatalf("unexpected token %q; want \"\\U00011F02\"", token)
	}
}

func formatBytes(buf []byte) string {
	strs := make([]string, len(buf))
	for i, b := range buf {
		strs[i] = fmt.Sprintf("0x%02x", b)
	}
	return strings.Join(strs, "   ")
}

func formatByteRanges(bufs [][]byte) string {
	strs := make([]string, len(bufs))
	for i, b := range bufs {
		strs[i] = formatBytes(b)
	}
	return strings.Join(strs, " | ")
}

// dropReader is a reader that returns at most one byte per call, thereby
// acting as a worst-case scenario for handling streaming data over a socket.
type dripReader struct {
	remain []byte
}

func (r *dripReader) Read(buf []byte) (int, error) {
	if len(r.remain) == 0 {
		return 0, io.EOF
	}
	if len(buf) == 0 {
		return 0, nil
	}
	buf[0] = r.remain[0]
	r.remain = r.remain[1:]
	return 1, nil
}
