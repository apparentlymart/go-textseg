package charprops

import "fmt"

// CharProperties is effectively a tuple of both a [GCBProperty] and an
// [InCBProperty] value, stored compactly as a bitwise-OR of valid values
// of those two types.
type CharProperties uint8

func MakeCharProperties(gcpProp GCBProperty, emojiProp GCBProperty, inCBProp InCBProperty) CharProperties {
	return CharProperties(uint8(gcpProp) | uint8(emojiProp) | uint8(inCBProp))
}

func (cp CharProperties) GCBProperty() GCBProperty {
	return GCBProperty(cp & 0x0f)
}

func (cp CharProperties) InCBProperty() InCBProperty {
	return InCBProperty(cp & 0xf0)
}

// GCBProperty is an enumeration of Grapheme_Cluster_Break property values,
// from [UAX#29 Section 3.1].
//
// The ExtendedPictographic property is actually derived from the Emoji
// standard's character tables, but is treated by UAX#29 as mutually-exclusive
// with the grapheme cluster break property value and so included in this
// enumeration for simplicity's sake.
//
// Note that the values of this type only set bits in the low nybble because
// they are intended to be bitwise-ORed with GCBProperty to produce
// [CharProperties] values.
//
// [UAX#29 Section 3.1]: https://www.unicode.org/reports/tr29/#Grapheme_Cluster_Break_Property_Values
type GCBProperty uint8

const (
	GCBNone                 GCBProperty = 0x00
	GCBCR                   GCBProperty = 0x01
	GCBControl              GCBProperty = 0x02
	GCBExtend               GCBProperty = 0x03
	GCBExtendedPictographic GCBProperty = 0x04
	GCBL                    GCBProperty = 0x05
	GCBLF                   GCBProperty = 0x06
	GCBLV                   GCBProperty = 0x07
	GCBLVT                  GCBProperty = 0x08
	GCBPrepend              GCBProperty = 0x09
	GCBRegionalIndicator    GCBProperty = 0x0a
	GCBSpacingMark          GCBProperty = 0x0b
	GCBT                    GCBProperty = 0x0c
	GCBV                    GCBProperty = 0x0d
	GCBZWJ                  GCBProperty = 0x0e
)

// LookupGCBProperty returns the [GCPProperty] value corresponding to the
// given property name, or [GCBNone] if the name is not recognized.
func LookupGCBProperty(name string) GCBProperty {
	switch name {
	case "CR":
		return GCBCR
	case "Control":
		return GCBControl
	case "Extend":
		return GCBExtend
	case "L":
		return GCBL
	case "LF":
		return GCBLF
	case "LV":
		return GCBLV
	case "LVT":
		return GCBLVT
	case "Prepend":
		return GCBPrepend
	case "Regional_Indicator":
		return GCBRegionalIndicator
	case "SpacingMark":
		return GCBSpacingMark
	case "T":
		return GCBT
	case "V":
		return GCBV
	case "ZWJ":
		return GCBZWJ
	default:
		return GCBNone
	}
}

// LookupEmojiProperty returns the [GCPProperty] value corresponding to the
// given emoji property name, or [GCBNone] if the name is not recognized.
//
// Note that because the segmentation rules treat the relevant emoji properties
// as mutually exclusive with the main grapheme clustering properties, we
// use [GCPProperty] to represent both of them and assume that no character
// will have both a nonzero result from [LookupGCPProperty] and from
// [LookupEmojiProperty], and so it's okay to collapse them together into the
// same enumeration. If they ever stop being mutually-exclusive then we'll
// need a different strategy.
func LookupEmojiProperty(name string) GCBProperty {
	switch name {
	case "Extended_Pictographic":
		return GCBExtendedPictographic
	default:
		return GCBNone // we don't need any other properties for our purposes here
	}
}

func (p GCBProperty) IsAnyControl() bool {
	switch p {
	case GCBLF, GCBCR, GCBControl:
		return true
	default:
		return false
	}
}

func (p GCBProperty) String() string {
	switch p {
	case GCBCR:
		return "CR"
	case GCBControl:
		return "Control"
	case GCBExtend:
		return "Extend"
	case GCBL:
		return "L"
	case GCBLF:
		return "LF"
	case GCBLV:
		return "LV"
	case GCBLVT:
		return "LVT"
	case GCBPrepend:
		return "Prepend"
	case GCBRegionalIndicator:
		return "Regional_Indicator"
	case GCBSpacingMark:
		return "SpacingMark"
	case GCBT:
		return "T"
	case GCBV:
		return "V"
	case GCBZWJ:
		return "ZWJ"
	case GCBExtendedPictographic:
		return "Extended_Pictographic"
	case GCBNone:
		return "None"
	default:
		return fmt.Sprintf("0x%02x", uint8(p))
	}
}

// InCBProperty is an enumeration of Indic_Conjunct_Break property values,
// as defined in DerivedCoreProperties.txt based on
// [the rules in UAX#44].
//
// These are used in the rule that avoids splitting orthographic syllables in
// inappropriated ways, [GB9c].
//
// Note that the values of this type only set bits in the high nybble because
// they are intended to be bitwise-ORed with GCBProperty to produce
// [CharProperties] values.
//
// [the rules in UAX#44]: https://www.unicode.org/reports/tr44/#Indic_Conjunct_Break
// [GB9c]: https://www.unicode.org/reports/tr29/#GB9c
type InCBProperty uint8

const (
	InCBNone      InCBProperty = 0x00
	InCBConsonant InCBProperty = 0x10
	InCBExtend    InCBProperty = 0x20
	InCBLinker    InCBProperty = 0x30
)

// LookupGCBProperty returns the [InCBProperty] value corresponding to the
// given property name, or [InCBNone] if the name is not recognized.
func LookupInCBProperty(name string) InCBProperty {
	switch name {
	case "Consonant":
		return InCBConsonant
	case "Extend":
		return InCBExtend
	case "Linker":
		return InCBLinker
	default:
		return InCBNone
	}
}
