package gentable

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"unicode/utf8"

	"github.com/apparentlymart/go-textseg/v16/textseg/internal/charprops"
	"github.com/apparentlymart/go-textseg/v16/textseg/internal/charprops/ucdparse"
)

const blockShift = 6
const blockSize = 1 << blockShift
const blockMask = uint8(blockSize - 1)

var zeroIndexBlock = make([]uint16, blockSize)
var zeroPropertyBlock = make([]uint8, blockSize)

func BuildRawPropertyTree(baseDir string) (RawPropertyTree, error) {
	gbpPath := filepath.Join(baseDir, "auxiliary", "GraphemeBreakProperty.txt")
	emojiPath := filepath.Join(baseDir, "emoji", "emoji-data.txt")
	dcpPath := filepath.Join(baseDir, "DerivedCoreProperties.txt")
	ret := RawPropertyTree{}

	gbpFile, err := os.Open(gbpPath)
	if err != nil {
		return ret, err
	}
	emojiFile, err := os.Open(emojiPath)
	if err != nil {
		return ret, err
	}
	dcpFile, err := os.Open(dcpPath)
	if err != nil {
		return ret, err
	}

	tip := newTreeInProgress()
	gpbSc := ucdparse.NewScanner(gbpFile)
	emojiSc := ucdparse.NewScanner(emojiFile)
	dcpSc := ucdparse.NewScanner(dcpFile)
	for {
		entry, err := gpbSc.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			return ret, fmt.Errorf("reading GraphemeBreakProperty.txt: %w", err)
		}
		propName := entry.FirstField()
		prop := charprops.LookupGCBProperty(propName)
		for r := entry.Start; r <= entry.End; r++ {
			tip.BitwiseOr(r, uint8(prop))
		}
	}
	for {
		entry, err := emojiSc.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			return ret, fmt.Errorf("reading emoji-data.txt: %w", err)
		}
		propName := entry.FirstField()
		prop := charprops.LookupEmojiProperty(propName)
		if prop == 0 {
			continue
		}
		for r := entry.Start; r <= entry.End; r++ {
			tip.BitwiseOr(r, uint8(prop))
		}
	}
	for {
		entry, err := dcpSc.NextEntry()
		if err == io.EOF {
			break
		} else if err != nil {
			return ret, fmt.Errorf("reading DerivedCoreProperties.txt: %w", err)
		}
		fields := slices.Collect(entry.AllFields())
		if len(fields) < 2 || fields[0] != "InCB" {
			continue // irrelevant to our goals, then
		}
		propName := fields[1]
		prop := charprops.LookupInCBProperty(propName)
		for r := entry.Start; r <= entry.End; r++ {
			tip.BitwiseOr(r, uint8(prop))
		}
	}

	ret = tip.Compact()
	return ret, nil
}

type treeInProgress struct {
	rawProps []uint8
	indices  []uint16
}

func newTreeInProgress() *treeInProgress {
	return &treeInProgress{
		// The first 128 property slots are reserved for direct lookup of the
		// ASCII characters, so we'll preallocate the storage for that.
		rawProps: make([]uint8, 128),

		// The first 64 index slots are reserved for the indices of the
		// remaining valid initial bytes, which are indices into either
		// rawProps directly or to a block offset elsewhere in indices
		// depending on whether it's a two-byte or longer UTF-8 encoding.
		indices: make([]uint16, blockSize),
	}
}

func (t *treeInProgress) BitwiseOr(r rune, propMask uint8) {
	rawProps := t.ensureProps(r)
	*rawProps |= propMask
}

func (t *treeInProgress) ensureProps(r rune) *uint8 {
	asU8 := make([]byte, 4)
	byteLen := utf8.EncodeRune(asU8, r)
	asU8 = asU8[:byteLen]
	if byteLen == 1 {
		// The ASCII characters are treated as a special case directly looked
		// up in the first 128 elements of rawProps, which are immediately
		// allocated at the creation of a [treeInProgress].
		return &t.rawProps[int(asU8[0])]
	}

	// In all of the following we just discard the top two bits of each
	// byte, because in the first byte they are always 0b11 set and in
	// any continuation bytes they are always 0b10. That gives us a natural
	// block size of 64 entries, using the low 6 bits of each byte.

	nested := asU8[:len(asU8)-2]
	penultimate := asU8[len(asU8)-2]
	final := asU8[len(asU8)-1]
	currentBlockIdx := 0
	for _, b := range nested {
		indexBlockIndexIdx := (currentBlockIdx << blockShift) + int(b&blockMask)
		indexBlockIndex := int(t.indices[indexBlockIndexIdx])
		if indexBlockIndex == 0 {
			// Need to allocate a new index block, then.
			indexBlockIndex = len(t.indices) >> blockShift
			t.indices = append(t.indices, zeroIndexBlock...)
			t.indices[indexBlockIndexIdx] = uint16(indexBlockIndex)
		}
		currentBlockIdx = indexBlockIndex
	}
	var propBlockIndex int
	{
		propBlockIndexIdx := (currentBlockIdx << blockShift) + int(penultimate&blockMask)
		propBlockIndex = int(t.indices[propBlockIndexIdx])
		if propBlockIndex == 0 {
			// Need to allocate a new property block, then.
			propBlockIndex = len(t.rawProps) >> blockShift
			t.rawProps = append(t.rawProps, zeroPropertyBlock...)
			t.indices[propBlockIndexIdx] = uint16(propBlockIndex)
		}
	}
	return &t.rawProps[(propBlockIndex<<blockShift)+int(final&blockMask)]
}

func (t *treeInProgress) Compact() RawPropertyTree {
	var props []uint8
	var indices []uint16
	propBlockIdxs := make(map[string]uint16)
	indexBlockIdxs := make(map[string]uint16)

	// The first 128 property values are fixed in place to represent the
	// properties for the ASCII characters.
	propBlockIdxs[propertyBlockKey(t.rawProps[:blockSize])] = 0
	propBlockIdxs[propertyBlockKey(t.rawProps[blockSize:blockSize*2])] = 1
	props = append(props, t.rawProps[:128]...)

	// The first block of indices is the root, but we're going to rebuild
	// that here as we reassign the block ids during compaction, so we'll
	// start out with it all zeroed and then fill out in as we work.
	indices = make([]uint16, blockSize)
	dstRootIndices := indices // the first block is the root, which will fill out as we go
	srcRootIndices := t.indices[:blockSize]

	ret := RawPropertyTree{
		Indices:    indices,
		Properties: props,
	}

	// With those fixed elements in place, our task now is to walk the non-ASCII
	// part of the tree, visiting all of the blocks, and rebuild the same tree
	// while consolidating any duplicate blocks.
	for i := range srcRootIndices {
		var indexLevels int
		// "i" here is a UTF-8 initial byte with its top two bits zeroed out,
		// so the following is just the UTF-8 length rules but adapted for
		// that masking.
		switch i >> 4 {
		case 0b00, 0b01: // two bytes long
			indexLevels = 0
		case 0b10: // three bytes long
			indexLevels = 1
		case 0b011: // four bytes long
			indexLevels = 2
		default:
			panic("invalid initial index")
		}
		t.buildCompactedBlocks(&ret, indexLevels, srcRootIndices, dstRootIndices, i, indexBlockIdxs, propBlockIdxs)
	}

	return ret
}

func (t *treeInProgress) buildCompactedBlocks(ret *RawPropertyTree, indexLevels int, srcIndices, dstIndices []uint16, currentIdx int, propBlockIdxs, indexBlockIdxs map[string]uint16) {
	srcBlockIdx := srcIndices[currentIdx]
	if indexLevels > 0 {
		// We're dealing with an intermediate index block, whose elements
		// refer to child index blocks.

		var srcBlock []uint16
		if srcBlockIdx == 0 {
			// This is a subtree that has no properties at all then, and so
			// no block was allocated for it during initial construction.
			// We'll just make sure this eventually points to an all-zeroes
			// property block.
			srcBlock = make([]uint16, blockSize)
		} else {
			srcBlockStart := srcBlockIdx << blockShift
			srcBlock = t.indices[srcBlockStart : srcBlockStart+blockSize]
		}

		newBlock := make([]uint16, blockSize)
		// We'll now recursively populate the new block before we decide
		// whether it's unique or if we have an existing block that already
		// matches it.
		for i := range srcBlock {
			t.buildCompactedBlocks(ret, indexLevels-1, srcBlock, newBlock, i, indexBlockIdxs, propBlockIdxs)
		}
		blockKey := indexBlockKey(newBlock)
		newBlockIndex, ok := indexBlockIdxs[blockKey]
		if !ok {
			// This block is different to all others we saw before, so we'll
			// add it and record its index for potential reuse later.
			newBlockIndex = uint16(len(ret.Indices) >> blockShift)
			ret.Indices = append(ret.Indices, newBlock...)
			indexBlockIdxs[blockKey] = newBlockIndex
		}
		dstIndices[currentIdx] = newBlockIndex
	}

	// We're dealing with a leaf index block, whose elements refer to property
	// blocks.
	var srcBlock []uint8
	if srcBlockIdx == 0 {
		// This is a subtree that has no properties at all then, and so
		// no block was allocated for it during initial construction.
		// We'll just make sure this eventually points to an all-zeroes
		// property block.
		srcBlock = make([]uint8, blockSize)
	} else {
		srcBlockStart := srcBlockIdx << blockShift
		srcBlock = t.rawProps[srcBlockStart : srcBlockStart+blockSize]
	}

	blockKey := propertyBlockKey(srcBlock)
	newBlockIndex, ok := propBlockIdxs[blockKey]
	if !ok {
		// This block is different to all others we saw before, so we'll
		// add it and record its index for potential reuse later.
		newBlockIndex = uint16(len(ret.Properties) >> blockShift)
		ret.Properties = append(ret.Properties, srcBlock...)
		propBlockIdxs[blockKey] = newBlockIndex
	}
	dstIndices[currentIdx] = newBlockIndex
}

func indexBlockKey(block []uint16) string {
	if len(block) != blockSize {
		panic("incorrect index block size")
	}
	var bytes [blockSize * 2]byte
	for i, v := range block {
		bytes[i*2] = byte(v)
		bytes[i*2+1] = byte(v >> 8)
	}
	return string(bytes[:])
}

func propertyBlockKey(block []uint8) string {
	if len(block) != blockSize {
		panic("incorrect property block size")
	}
	return string(block)
}
