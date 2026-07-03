package charprops

import (
	"unicode/utf8"
)

//go:generate go run ./generate.go

var zeroProperties CharProperties

func LookupFirstChar(p []byte) (props CharProperties, length int) {
	if len(p) == 0 {
		return zeroProperties, 0
	}

	first := p[0]
	if first < 128 {
		length = 1
	} else if (first & 0b11100000) == 0b11000000 {
		length = 2
	} else if (first & 0b11110000) == 0b11100000 {
		length = 3
	} else if (first & 0b11111000) == 0b11110000 {
		length = 4
	}

	switch length {
	case 1:
		return lookupProps[first], length
	case 2:
		blockIdx := int(lookupIndices[first&0b111111])
		return lookupProps[(blockIdx<<6)+int(p[1]&0b111111)], length
	case 3:
		blockIdx := int(lookupIndices[first&0b111111])
		blockIdx = int(lookupIndices[(blockIdx<<6)+int(p[1]&0b111111)])
		return lookupProps[(blockIdx<<6)+int(p[2]&0b111111)], length
	case 4:
		blockIdx := int(lookupIndices[first&0b111111])
		blockIdx = int(lookupIndices[(blockIdx<<6)+int(p[1]&0b111111)])
		blockIdx = int(lookupIndices[(blockIdx<<6)+int(p[2]&0b111111)])
		return lookupProps[(blockIdx<<6)+int(p[3]&0b111111)], length
	default:
		return zeroProperties, 0
	}
}

var runeErrorBytes []byte

func init() {
	re := make([]byte, 3)
	n := utf8.EncodeRune(re, utf8.RuneError)
	runeErrorBytes = re[:n]
}
