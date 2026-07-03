package charprops

import (
	"io"
	"os"
	"slices"
	"testing"
	"unicode/utf8"

	"github.com/apparentlymart/go-textseg/v16/textseg/internal/charprops/ucdparse"
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
			CharProperties(uint8(GCBExtend) | uint8(InCBExtend)),
			3,
		},
		{
			"\U0001F1E6\U0001F1E6",
			CharProperties(GCBRegionalIndicator),
			4,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			gotProperties, gotLength := LookupFirstChar([]byte(test.Input))
			if got, want := gotProperties, test.WantProperties; got != want {
				t.Errorf("wrong properties\ngot:  0x%02x\nwant: 0x%02x", got, want)
			}
			if got, want := gotLength, test.WantLength; got != want {
				t.Errorf("wrong length\ngot:  %#v\nwant: %#v", got, want)
			}
		})
	}
}

func TestLookupFirstChar_graphemeBreakProperty(t *testing.T) {
	// This test verifies that our generated property tree agrees with the
	// properties specified in GraphemeBreakProperty.txt, which should always
	// be true if the generator is correct because the generator uses this
	// same file as part of its input.
	f, err := os.Open("ucd/auxiliary/GraphemeBreakProperty.txt")
	if err != nil {
		t.Error(err)
	}

	sc := ucdparse.NewScanner(f)
	for {
		entry, err := sc.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		wantGBP := LookupGCBProperty(entry.FirstField())
		for r := entry.Start; r < entry.End; r++ {
			buf := make([]byte, 4)
			n := utf8.EncodeRune(buf, r)
			buf = buf[:n]
			t.Run(string(buf), func(t *testing.T) {
				t.Logf("for scalar value U+%04X", r)
				gotProperties, gotLength := LookupFirstChar(buf)
				if got, want := gotProperties.GCBProperty(), wantGBP; got != want {
					t.Errorf("wrong properties\ngot:  0x%02x\nwant: 0x%02x", got, want)
				}
				if got, want := gotLength, n; got != want {
					t.Errorf("wrong length\ngot:  %#v\nwant: %#v", got, want)
				}
			})
		}
	}
}

func TestLookupFirstChar_emojiExtendedPictographic(t *testing.T) {
	// This test verifies that our generated property tree agrees with the
	// Extended_Pictographic property specified in emoji-data.txt, which should
	// always be true if the generator is correct because the generator uses
	// this same file as part of its input.
	f, err := os.Open("ucd/emoji/emoji-data.txt")
	if err != nil {
		t.Error(err)
	}

	sc := ucdparse.NewScanner(f)
	for {
		entry, err := sc.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		wantGBP := LookupEmojiProperty(entry.FirstField())
		if wantGBP != GCBExtendedPictographic {
			continue // other emoji properties are irrelevant
		}
		for r := entry.Start; r < entry.End; r++ {
			buf := make([]byte, 4)
			n := utf8.EncodeRune(buf, r)
			buf = buf[:n]
			t.Run(string(buf), func(t *testing.T) {
				t.Logf("for scalar value U+%04X", r)
				gotProperties, gotLength := LookupFirstChar(buf)
				if got, want := gotProperties.GCBProperty(), wantGBP; got != want {
					t.Errorf("wrong properties\ngot:  0x%02x\nwant: 0x%02x", got, want)
				}
				if got, want := gotLength, n; got != want {
					t.Errorf("wrong length\ngot:  %#v\nwant: %#v", got, want)
				}
			})
		}
	}
}

func TestLookupFirstChar_indicConjunctBreakProperty(t *testing.T) {
	// This test verifies that our generated property tree agrees with the
	// InCB property values specified in DerivedCoreProperties.txt, which should
	// always be true if the generator is correct because the generator uses
	// this same file as part of its input.
	f, err := os.Open("ucd/DerivedCoreProperties.txt")
	if err != nil {
		t.Error(err)
	}

	sc := ucdparse.NewScanner(f)
	for {
		entry, err := sc.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		fields := slices.Collect(entry.AllFields())
		if len(fields) < 2 || fields[0] != "InCB" {
			continue
		}

		wantInCBP := LookupInCBProperty(fields[1])
		for r := entry.Start; r < entry.End; r++ {
			buf := make([]byte, 4)
			n := utf8.EncodeRune(buf, r)
			buf = buf[:n]
			t.Run(string(buf), func(t *testing.T) {
				t.Logf("for scalar value U+%04X", r)
				gotProperties, gotLength := LookupFirstChar(buf)
				if got, want := gotProperties.InCBProperty(), wantInCBP; got != want {
					t.Errorf("wrong properties\ngot:  0x%02x\nwant: 0x%02x", got, want)
				}
				if got, want := gotLength, n; got != want {
					t.Errorf("wrong length\ngot:  %#v\nwant: %#v", got, want)
				}
			})
		}
	}
}
