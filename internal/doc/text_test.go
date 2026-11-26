package doc

import (
	"fmt"
	"slices"
	"testing"
)

func TestTextSpans(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []piece
	}{
		{"empty document", "", nil},
		{"plain text", "hola", []piece{p("hola", Plain)}},

		{"heading", "# Hola", []piece{p("Hola", Heading)}},
		{"bold", "**hola**", []piece{p("hola", Bold)}},
		{"italic", "*hola*", []piece{p("hola", Italic)}},
		{"code", "`hola`", []piece{p("hola", Code)}},
		{"quote", "> hola", []piece{p("hola", Quote)}},

		// styles combine while going deeper into the tree...
		{"bold inside heading", "# Hola **mundo**", []piece{
			p("Hola ", Heading), p("mundo", Heading|Bold),
		}},
		{"bold and italic", "***hola***", []piece{p("hola", Bold|Italic)}},
		{"italic inside quote", "> *hola*", []piece{p("hola", Quote|Italic)}},

		// ...and go back when leaving a node.
		{"style ends after bold", "a **b** c", []piece{
			p("a ", Plain), p("b", Bold), p(" c", Plain),
		}},
		{"heading ends at the next block", "# A\n\nb", []piece{
			p("A", Heading), p("b", Plain),
		}},

		// goldmark may split one text into several nodes: keep them all
		{"star inside bold", "**a*b**", []piece{p("a*", Bold), p("b", Bold)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pieces(tt.src, textSpans([]byte(tt.src)))
			if !slices.Equal(got, tt.want) {
				t.Errorf("textSpans(%q)\n got  %v\n want %v", tt.src, got, tt.want)
			}
		})
	}
}

// HELPERS
// piece is a readable version of a Span: its text instead of its positions.
type piece struct {
	Text  string
	Style Style
}

func (p piece) String() string {
	return fmt.Sprintf("{%q %s}", p.Text, p.Style)
}

// pieces turns spans into pieces, so test tables can be written with text.
func pieces(src string, spans []Span) []piece {
	var out []piece
	for _, s := range spans {
		out = append(out, piece{src[s.Start:s.End], s.Style})
	}
	return out
}

func p(text string, style Style) piece { return piece{text, style} }
