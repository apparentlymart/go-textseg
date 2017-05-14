package textseg

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

var errInvalid = errors.New("invalid utf-8 encoding")

// ScanGraphemeClusters is a split function for bufio.Scanner that splits
// on grapheme cluster boundaries.
func ScanGraphemeClusters(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	fmt.Printf("ScanGraphemeClusters(%#v, %#v)\n", data, atEOF)

	i := 0
	state := 0
Rune:
	for {
		if i >= len(data) {
			break Rune
		}
		chr, l := utf8.DecodeRune(data[i:])
		if chr == utf8.RuneError && i > 0 {
			// Always break before invalid UTF-8 sequences
			break Rune
		}

		rt := _GraphemeRuneType(chr)
		remain := data[i+l:]
		var nextRt _GraphemeRuneRange
		if len(remain) > 0 {
			nextChr, _ := utf8.DecodeRune(remain)
			nextRt = _GraphemeRuneType(nextChr)
		}

		fmt.Printf("Processing 0x%04x (%d) in state %d\n", chr, i, state)

		switch state {
		case 0:
			switch rt {
			case _GraphemeCR: // GB3
				if i == 0 {
					state = 3
				} else {
					break Rune
				}
			case _GraphemeControl, _GraphemeLF: // GB4 & GB5
				if i == 0 {
					i += l
				}
				break Rune
			case _GraphemeL: // GB6
				state = 6
			case _GraphemeLV, _GraphemeV: // GB7
				state = 7
			case _GraphemeLVT, _GraphemeT: // GB8
				state = 8
			case _GraphemePrepend: // GB9b
				state = 0
			case _GraphemeE_Base, _GraphemeE_Base_GAZ: // GB10
				state = 10
			case _GraphemeZWJ: // GB11
				state = 11
			case _GraphemeRegional_Indicator: // GB12 & GB13
				state = 12
			default: // GB999
				switch nextRt {
				case _GraphemeExtend, _GraphemeZWJ: // GB9
					state = 0
				case _GraphemeSpacingMark: // GB9a
					state = 0
				default: // includes nextRt == nil when !hasNext
					i += l
					break Rune
				}
			}

		case 3: // GB3
			// Previous character was CR. If current character is LF then
			// we'll consume it, but otherwise we're at a cluster boundary
			// due to GB4.
			//
			// Either way we break, because we always break after LF or CR LF.
			if rt == _GraphemeLF {
				i += l
			}
			break Rune

		case 6: // GB6
			switch rt {
			case _GraphemeL:
				state = 6
			case _GraphemeV, _GraphemeLV:
				state = 7
			case _GraphemeLVT:
				state = 8
			default:
				// Break without consuming the current character
				break Rune
			}

		case 7: // GB7
			switch rt {
			case _GraphemeV:
				state = 7
			case _GraphemeT:
				state = 8
			default:
				// Break without consuming the current character
				break Rune
			}

		case 8: // GB8
			switch rt {
			case _GraphemeT:
				state = 8
			default:
				// Break without consuming the current character
				break Rune
			}

		case 10: // GB10
			switch rt {
			case _GraphemeExtend:
				state = 10
			case _GraphemeE_Modifier:
				i += l
				break Rune
			default:
				// Break without consuming the current character
				break Rune
			}

		case 11: // GB11
			switch rt {
			case _GraphemeGlue_After_Zwj:
				i += l
				break Rune
			case _GraphemeE_Base_GAZ:
				state = 10
			default:
				// Break without consuming the current character
				break Rune
			}

		case 12: // GB12
			switch rt {
			case _GraphemeRegional_Indicator:
				i += l
				break Rune
			default:
				// Break without consuming the current character
				break Rune
			}

		default:
			// We should never get here. If we do, we'll switch back to state
			// 0 without advancing and retry.
			state = 0
			continue Rune
		}

		i += l
	}

	fmt.Printf("ScanGraphemeClusters result %#v, %#v, nil\n", i, data[:i])
	return i, data[:i], nil
}
