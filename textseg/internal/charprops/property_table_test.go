package charprops

import (
	"testing"
)

func TestLookupFirstChar(t *testing.T) {
	tests := []struct {
		Input          string
		WantProperties CharProperties
		WantLength     int
	}{
		{
			"",
			CharProperties(0),
			0,
		},
		{
			"\n",
			CharProperties(GCBLF),
			1,
		},
		{
			"\u2665\uFE0F",
			CharProperties(GCBExtendedPictographic),
			3,
		},
		{
			"\uFE0F",
			CharProperties(GCBControl),
			3,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			gotProperties, gotLength := LookupFirstChar([]byte(test.Input))
			if got, want := gotProperties, test.WantProperties; got != want {
				t.Errorf("wrong properties\ngot:  %#v\nwant: %#v", got, want)
			}
			if got, want := gotLength, test.WantLength; got != want {
				t.Errorf("wrong length\ngot:  %#v\nwant: %#v", got, want)
			}
		})
	}
}
