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
