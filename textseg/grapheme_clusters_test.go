package textseg

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestScanGraphemeClusters(t *testing.T) {
	tests := unicodeGraphemeTests

	for i, test := range tests {
		t.Run(fmt.Sprintf("%03d-%x", i, test.input), func(t *testing.T) {
			got, err := AllTokens(test.input, ScanGraphemeClusters)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !reflect.DeepEqual(got, test.output) {
				// Also get the rune values resulting from decoding utf8,
				// since they are generally easier to look up to figure out
				// what's failing.
				runes := make([]string, 0, len(test.input))
				buf := test.input
				for len(buf) > 0 {
					r, size := utf8.DecodeRune(buf)
					runes = append(runes, fmt.Sprintf("0x%04x", r))
					buf = buf[size:]
				}

				t.Errorf(
					"wrong result\ninput: %#v\nrunes: %s\ngot:   %#v\nwant:  %#v",
					test.input, strings.Join(runes, " "), got, test.output,
				)
			}
		})
	}
}
