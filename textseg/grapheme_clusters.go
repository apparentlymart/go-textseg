package textseg

import (
	"errors"

	"github.com/apparentlymart/go-textseg/v16/textseg/internal/charprops"
	"github.com/apparentlymart/go-textseg/v16/textseg/internal/machine"
)

// Error is a legacy symbol that is no longer used.
//
// Deprecated: This will be removed in a future release.
var Error = errors.New("invalid UTF8 text")

// ScanGraphemeClusters is a split function for bufio.Scanner that splits
// on grapheme cluster boundaries.
func ScanGraphemeClusters(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	advance, token, err := ScanUTF8Sequences(data, atEOF)
	if err != nil || (advance == 0 && len(token) == 0) {
		return advance, token, err
	}

	properties, count := charprops.LookupFirstChar(token)
	if count != advance {
		// An invalid UTF-8 sequence then, so we'll just report the
		// next byte standalone.
		return 1, data[0:1], nil
	}
	remain := data[advance:]

	state := machine.StateBase
	prev := properties
	for {
		if len(remain) == 0 {
			if atEOF {
				// If we're at the end of the file then whatever we've
				// accumulated so far is a grapheme cluster.
				return count, data[:count], nil
			}
			// If we're not at the end of the file then we'll need more
			// bytes before we can decide if we've reached the end of a
			// grapheme cluster.
			return 0, nil, nil
		}

		advance, token, err := ScanUTF8Sequences(remain, atEOF)
		if err != nil {
			// If the next sequence is incomplete or invalid then we'll
			// just return here
			return count, data[:count], err
		} else if advance == 0 && len(token) == 0 {
			// More bytes required to complete the next UTF-8 sequence.
			return 0, nil, nil
		}

		next, moreCount := charprops.LookupFirstChar(token)
		if moreCount != advance {
			// An invalid UTF-8 sequence then, so we'll just report what
			// we found so far and let the next round deal with the invalid
			// prefix.
			return count, data[:count], nil
		}
		remain = remain[advance:]

		split, nextState := state.Transition(prev, next)
		if split {
			// We've found the next split point, so we'll just return what
			// we have and then let the next call pick up from here.
			return count, data[:count], nil
		}
		count += moreCount
		state = nextState
		prev = next
	}
}
