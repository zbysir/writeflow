package llm

import (
	"github.com/tmc/langchaingo/textsplitter"
)

type Spliter interface {
	Split(s string) (ss []string)
}

type MarkDoneSplit struct {
	size int
}

func NewMarkDoneSplit(size int) *MarkDoneSplit {
	if size == 0 {
		panic("size can't be zero")
	}
	return &MarkDoneSplit{size: size}
}

func (d *MarkDoneSplit) Split(s string) (ss []string) {
	// https://github.com/qalqi/langchainjs/blob/17bfe9b0a2468c813042ad5888b4eb72294942ca/langchain/src/text_splitter.ts#L310
	sp := textsplitter.RecursiveCharacter{
		Separators: []string{
			"\n## ",
			"\n### ",
			"\n#### ",
			"\n##### ",
			"\n###### ",
			// Note the alternative syntax for headings (below) is not handled here
			// Heading level 2
			// ---------------
			// End of code block
			"```\n\n",
			// Horizontal lines
			"\n\n***\n\n",
			"\n\n---\n\n",
			"\n\n___\n\n",
			// Note that this splitter doesn't handle horizontal lines defined
			// by *three or more* of ***, ---, or ___, but this is not handled
			"\n\n",
			"\n",
			" ",
			""},
		ChunkSize:    d.size,
		ChunkOverlap: 200,
	}
	ss, err := sp.SplitText(s)
	if err != nil {
		panic(err)
	}

	return ss
}
