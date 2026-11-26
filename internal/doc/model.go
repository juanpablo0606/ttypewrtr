package doc

import "strings"

// Kind is what a paragraph is.
type Kind int

const (
	KindBody     Kind = iota // normal text
	KindHeading              // a heading; its level is in Para.Level
	KindQuote                // a quote
	KindListItem             // an item of a list
)

// Align is how a paragraph is aligned.
type Align int

const (
	AlignDefault Align = iota // the document's default alignment
	AlignLeft
	AlignCenter
	AlignRight
	AlignJustify
)

// Para is one paragraph, which is one line in the editor.
type Para struct {
	Kind  Kind
	Level int // heading level (1–6); 0 for other kinds
	Align Align
	Text  string // the text the writer sees, without markdown symbols
	Spans []Span // bold, italic and code ranges inside Text
}

// Doc is a whole document: a list of paragraphs.
type Doc struct {
	Paras []Para
}

// Text returns the text of all paragraphs, one per line.
func (d Doc) Text() string {
	texts := make([]string, len(d.Paras))
	for i, p := range d.Paras {
		texts[i] = p.Text
	}
	return strings.Join(texts, "\n")
}
