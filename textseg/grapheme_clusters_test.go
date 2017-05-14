package textseg

import (
	"fmt"
	"reflect"
	"testing"
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
				t.Errorf(
					"wrong result\ninput: %#v\ngot:   %#v\nwant:  %#v",
					test.input, got, test.output,
				)
			}
		})
	}
}
