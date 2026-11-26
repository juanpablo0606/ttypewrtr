package render

import (
	"slices"
	"testing"
)

func TestSpans(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []piece
	}{
		{"empty document", "", nil},
		{"plain text", "hola mundo", []piece{p("hola mundo", Plain)}},

		// headings
		{"heading", "# Hola", []piece{
			p("# ", Marker), p("Hola", Heading),
		}},
		{"heading with bold", "# Hola **mundo**", []piece{
			p("# ", Marker), p("Hola ", Heading), p("**", Marker),
			p("mundo", Heading|Bold), p("**", Marker),
		}},

		// emphasis
		{"bold", "**hola**", []piece{
			p("**", Marker), p("hola", Bold), p("**", Marker),
		}},
		{"italic", "*hola*", []piece{
			p("*", Marker), p("hola", Italic), p("*", Marker),
		}},
		{"bold and italic: neighbor markers are one span", "***hola***", []piece{
			p("***", Marker), p("hola", Bold|Italic), p("***", Marker),
		}},
		{"bold in the middle", "a **b** c", []piece{
			p("a ", Plain), p("**", Marker), p("b", Bold), p("**", Marker), p(" c", Plain),
		}},

		// the cases that are hard to get right by hand
		{"star inside bold: one span", "**a*b**", []piece{
			p("**", Marker), p("a*b", Bold), p("**", Marker),
		}},
		{"escaped stars are not italic", `\*no\*`, []piece{
			p(`\*no\*`, Plain),
		}},
		{"no bold inside code", "`**x**`", []piece{
			p("`", Marker), p("**x**", Code), p("`", Marker),
		}},

		// blocks
		{"quote", "> hola", []piece{
			p("> ", Marker), p("hola", Quote),
		}},
		{"bold inside quote: neighbor markers are one span", "> **hola**", []piece{
			p("> **", Marker), p("hola", Quote|Bold), p("**", Marker),
		}},
		{"list item", "- hola", []piece{
			p("- ", Marker), p("hola", Plain),
		}},

		// what is between blocks (newlines) is a Marker too
		{"two paragraphs", "uno\n\ndos", []piece{
			p("uno", Plain), p("\n\n", Marker), p("dos", Plain),
		}},
		{"heading then paragraph", "# A\n\nb", []piece{
			p("# ", Marker), p("A", Heading), p("\n\n", Marker), p("b", Plain),
		}},
		{"newline at the end", "hola\n", []piece{
			p("hola", Plain), p("\n", Marker),
		}},

		// positions are bytes: accents must not be split
		{"accents", "**canción**", []piece{
			p("**", Marker), p("canción", Bold), p("**", Marker),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pieces(tt.src, Spans([]byte(tt.src)))
			if !slices.Equal(got, tt.want) {
				t.Errorf("Spans(%q)\n got  %v\n want %v", tt.src, got, tt.want)
			}
		})
	}
}

// The rules every result of Spans must follow, checked on many documents.
func TestSpansInvariants(t *testing.T) {
	docs := []string{
		"", "hola", "# Hola **mundo** y *vos*",
		"> cita **x**\n\n- item\n\n**a*b** \\*no\\* `**code**`\n",
		"uno\ndos\n\ntres", "***hola***", "# A\n\n## B\n\nc *d* e",
		"**canción** ñandú", "texto final sin salto",
		"texto con salto final\n",
	}

	for _, doc := range docs {
		spans := Spans([]byte(doc))
		pos := 0
		for i, s := range spans {
			// 1. They cover the whole document, in order, with no gaps or overlaps.
			if s.Start != pos {
				t.Errorf("Spans(%q)[%d] starts at %d, want %d (gap or overlap)", doc, i, s.Start, pos)
			}
			// 2. No empty spans.
			if s.End <= s.Start {
				t.Errorf("Spans(%q)[%d] is empty: [%d:%d]", doc, i, s.Start, s.End)
			}
			// 3. Two neighbors never have the same style: they should be one span.
			if i > 0 && spans[i-1].Style == s.Style {
				t.Errorf("Spans(%q)[%d] and [%d] have the same style %s: merge them", doc, i-1, i, s.Style)
			}
			pos = s.End
		}
		if pos != len(doc) {
			t.Errorf("Spans(%q) ends at %d, want %d (end of document)", doc, pos, len(doc))
		}
	}
}
