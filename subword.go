package hangulize

import (
	"bytes"
)

// Subword is a chunk of a word with a level number. The level indicates which
// pipeline step generated this sw.
type Subword struct {
	Word  string
	Level int
}

// Builder is a buffer to build a []Subword array.
type Builder struct {
	subwords []Subword
}

// String() concatenates buffered subwords to assemble the full word.
func (b *Builder) String() string {
	var buf bytes.Buffer
	for _, sw := range b.subwords {
		buf.WriteString(sw.Word)
	}
	return buf.String()
}

// Append extends the underlying subwords by the given ones.
func (b *Builder) Append(subwords ...Subword) {
	b.subwords = append(b.subwords, subwords...)
}

// Reset discards the underlying subwords.
func (b *Builder) Reset() {
	b.subwords = b.subwords[:0]
}

// Subwords builds the buffered subwords into a []Subword array. It merges
// adjoin subwords if they share the same level.
func (b *Builder) Subwords() []Subword {
	var subwords []Subword

	if len(b.subwords) == 0 {
		// No subwords buffered.
		return subwords
	}

	// Merge same level adjoin subwords.
	var buf bytes.Buffer
	mergingLevel := -1

	for _, sw := range b.subwords {
		if sw.Level != mergingLevel && mergingLevel != -1 {
			// Keep the merged sw.
			merged := &Subword{buf.String(), mergingLevel}
			subwords = append(subwords, *merged)

			// Open a new one.
			buf.Reset()
		}

		buf.WriteString(sw.Word)
		mergingLevel = sw.Level
	}

	merged := &Subword{buf.String(), mergingLevel}
	subwords = append(subwords, *merged)

	return subwords
}
