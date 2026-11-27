package doc

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Para
	}{
		{"empty file", "", nil},
		{"one paragraph", "hola", []Para{
			{Text: "hola"},
		}},
		{"two paragraphs", "uno\n\ndos", []Para{
			{Text: "uno"},
			{Text: "dos"},
		}},
		{"paragraph over two lines is one paragraph", "uno\ndos", []Para{
			{Text: "uno dos"},
		}},

		// headings
		{"heading", "# Título", []Para{
			{Kind: KindHeading, Level: 1, Text: "Título"},
		}},
		{"heading level 3", "### Sub", []Para{
			{Kind: KindHeading, Level: 3, Text: "Sub"},
		}},
		{"heading then text", "# Título\n\ntexto", []Para{
			{Kind: KindHeading, Level: 1, Text: "Título"},
			{Text: "texto"},
		}},

		// quotes: each paragraph inside the quote is a Quote paragraph
		{"quote", "> cita", []Para{
			{Kind: KindQuote, Text: "cita"},
		}},
		{"quote with two paragraphs", "> uno\n>\n> dos", []Para{
			{Kind: KindQuote, Text: "uno"},
			{Kind: KindQuote, Text: "dos"},
		}},

		// lists: each item is a ListItem paragraph
		{"list", "- a\n- b", []Para{
			{Kind: KindListItem, Text: "a"},
			{Kind: KindListItem, Text: "b"},
		}},
		{"list with blank lines between items", "- a\n\n- b", []Para{
			{Kind: KindListItem, Text: "a"},
			{Kind: KindListItem, Text: "b"},
		}},

		// styles inside paragraphs come from inline
		{"styles", "**hola** mundo", []Para{
			{Text: "hola mundo", Spans: []Span{{0, 4, Bold}}},
		}},
		{"styles in a heading", "# Hola *vos*", []Para{
			{Kind: KindHeading, Level: 1, Text: "Hola vos", Spans: []Span{{5, 8, Italic}}},
		}},
		{"styles in a list item", "- **a** b", []Para{
			{Kind: KindListItem, Text: "a b", Spans: []Span{{0, 1, Bold}}},
		}},

		// a small story
		{"mixed document", "# Capítulo 1\n\nLa *lluvia* no paraba.\n\n> Nadie vino.\n\n- uno\n- dos\n", []Para{
			{Kind: KindHeading, Level: 1, Text: "Capítulo 1"},
			{Text: "La lluvia no paraba.", Spans: []Span{{3, 9, Italic}}},
			{Kind: KindQuote, Text: "Nadie vino."},
			{Kind: KindListItem, Text: "uno"},
			{Kind: KindListItem, Text: "dos"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Load([]byte(tt.src))
			if !reflect.DeepEqual(got.Paras, tt.want) {
				t.Errorf("Load(%q)\n got  %+v\n want %+v", tt.src, got.Paras, tt.want)
			}
		})
	}
}
