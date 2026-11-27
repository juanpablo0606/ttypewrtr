package doc

import (
	"slices"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// firstBlock parses src and returns its first block (a paragraph, a
// heading...), which is what inline receives.
func firstBlock(src []byte) ast.Node {
	doc := goldmark.DefaultParser().Parse(text.NewReader(src))
	return doc.FirstChild()
}

func TestInline(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		wantText  string
		wantSpans []Span
	}{
		{"plain text", "hola", "hola", nil},

		// symbols disappear, styles become spans over the clean text
		{"bold", "**hola** mundo", "hola mundo", []Span{{0, 4, Bold}}},
		{"italic", "a *b* c", "a b c", []Span{{2, 3, Italic}}},
		{"bold and italic", "***x***", "x", []Span{{0, 1, Bold | Italic}}},
		{"code", "`a*b*`", "a*b*", []Span{{0, 4, Code}}},
		{"two styled words", "**a** y *b*", "a y b", []Span{{0, 1, Bold}, {4, 5, Italic}}},
		{"neighbors with the same style are one span", "**a*b**", "a*b", []Span{{0, 3, Bold}}},

		// escapes and entities become the real characters
		{"escaped stars", `\*no\*`, "*no*", nil},
		{"entity", "&amp; y **x**", "& y x", []Span{{4, 5, Bold}}},

		// a paragraph is one line: line breaks become spaces
		{"soft line break", "uno\ndos", "uno dos", nil},
		{"line break inside bold", "**uno\ndos**", "uno dos", []Span{{0, 7, Bold}}},

		// the block itself adds no style: a heading's text is not Heading
		{"heading", "# Título **x**", "Título x", []Span{{8, 9, Bold}}},
		{"quote", "> una *cita*", "una cita", []Span{{4, 8, Italic}}},

		// positions are bytes: ó and ñ take 2 bytes each
		{"accents", "canción **ñandú**", "canción ñandú", []Span{{9, 16, Bold}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := []byte(tt.src)
			block := firstBlock(src)
			// For a quote, the text is inside its first paragraph.
			if _, ok := block.(*ast.Blockquote); ok {
				block = block.FirstChild()
			}

			gotText, gotSpans := inline(block, src)
			if gotText != tt.wantText {
				t.Errorf("inline(%q) text = %q, want %q", tt.src, gotText, tt.wantText)
			}
			if !slices.Equal(gotSpans, tt.wantSpans) {
				t.Errorf("inline(%q) spans = %v, want %v", tt.src, gotSpans, tt.wantSpans)
			}
		})
	}
}

// Every span must point to a real, non-empty part of the text.
func TestInlineSpansAreInsideText(t *testing.T) {
	docs := []string{
		"hola", "**hola** mundo", "***x***", "**a** y *b*", `\*no\*`,
		"**uno\ndos**", "# Título **x**", "canción **ñandú**", "a **b** *c* `d`",
	}

	for _, doc := range docs {
		src := []byte(doc)
		text, spans := inline(firstBlock(src), src)
		for _, s := range spans {
			if s.Start < 0 || s.End > len(text) || s.End <= s.Start {
				t.Errorf("inline(%q) = %q with bad span %v", doc, text, s)
			}
			if s.Style == Plain {
				t.Errorf("inline(%q) has a Plain span %v: only styled ranges are kept", doc, s)
			}
		}
	}
}
