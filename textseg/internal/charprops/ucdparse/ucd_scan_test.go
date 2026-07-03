package ucdparse

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

func TestScanner(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "ExampleProperties.txt"))
	if err != nil {
		t.Fatal(err)
	}

	type Item struct {
		Start, End rune
		FirstField string
		AllFields  []string
	}

	sc := NewScanner(f)
	var got []Item
	for {
		entry, err := sc.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("unexpected error for item %d: %s", len(got), err)
		}
		got = append(got, Item{
			Start:      entry.Start,
			End:        entry.End,
			FirstField: entry.FirstField(),
			AllFields:  slices.Collect(entry.AllFields()),
		})
	}
	want := []Item{
		{
			Start: 0x0000, End: 0x0080,
			FirstField: "Basic Latin",
			AllFields:  []string{"Basic Latin"},
		},
		{
			Start: 0x0080, End: 0x0100,
			FirstField: "Latin-1 Supplement",
			AllFields:  []string{"Latin-1 Supplement"},
		},
		{
			Start: 0x2764, End: 0x2765,
			FirstField: "Heart",
			AllFields:  []string{"Heart"},
		},
		{
			Start: 0x1D100, End: 0x1D200,
			FirstField: "Music",
			AllFields:  []string{"Music", "Astral"},
		},
	}

	for i := range want {
		wantItem := want[i]
		var gotItem Item
		if i < len(got) {
			gotItem = got[i]
		}
		// Using reflect.DeepEqual here just because this codebase avoids
		// having any non-stdlib Go module dependencies.
		if !reflect.DeepEqual(wantItem, gotItem) {
			t.Errorf("wrong item %d\ngot:  %#v\nwant: %#v", i, gotItem, wantItem)
		}
	}
	if len(got) > len(want) {
		excess := got[len(want):]
		for i, gotItem := range excess {
			t.Errorf("unexpected extra item %d: %#v", len(want)+i, gotItem)
		}
	}
}
