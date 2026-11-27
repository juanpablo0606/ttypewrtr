package doc

import (
	"reflect"
	"testing"
)

func TestParseFence(t *testing.T) {
	tests := []struct {
		text    string
		want    Align
		isFence bool
	}{
		// opening fences
		{"::: left", AlignLeft, true},
		{"::: center", AlignCenter, true},
		{"::: right", AlignRight, true},
		{"::: justify", AlignJustify, true},

		// closing fence: back to the document's default
		{":::", AlignDefault, true},

		// not fences
		{"hola", AlignDefault, false},
		{"", AlignDefault, false},
		{"::: middle", AlignDefault, false},       // unknown alignment
		{":::center", AlignDefault, false},        // needs a space
		{"texto ::: center", AlignDefault, false}, // must be the whole line
		{"::", AlignDefault, false},
	}

	for _, tt := range tests {
		got, isFence := parseFence(tt.text)
		if got != tt.want || isFence != tt.isFence {
			t.Errorf("parseFence(%q) = (%d, %v), want (%d, %v)", tt.text, got, isFence, tt.want, tt.isFence)
		}
	}
}

func TestLoadAlignAndEmpty(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Para
	}{
		// alignment fences: the paragraphs inside get the alignment, and the
		// fence lines themselves are not paragraphs
		{"centered heading", "::: center\n\n# Capítulo 1\n\n:::\n\ntexto", []Para{
			{Kind: KindHeading, Level: 1, Align: AlignCenter, Text: "Capítulo 1"},
			{Text: "texto"},
		}},
		{"right-aligned paragraph", "uno\n\n::: right\n\n—Nadie.\n\n:::\n\ndos", []Para{
			{Text: "uno"},
			{Align: AlignRight, Text: "—Nadie."},
			{Text: "dos"},
		}},
		{"several paragraphs in one fence", "::: justify\n\nuno\n\ndos\n\n:::", []Para{
			{Align: AlignJustify, Text: "uno"},
			{Align: AlignJustify, Text: "dos"},
		}},
		{"two fences in a row", "::: center\n\na\n\n:::\n\n::: right\n\nb\n\n:::", []Para{
			{Align: AlignCenter, Text: "a"},
			{Align: AlignRight, Text: "b"},
		}},
		{"fence around a quote", "::: center\n\n> cita\n\n:::", []Para{
			{Kind: KindQuote, Align: AlignCenter, Text: "cita"},
		}},
		{"fence never closed: lasts until the end", "::: center\n\na\n\nb", []Para{
			{Align: AlignCenter, Text: "a"},
			{Align: AlignCenter, Text: "b"},
		}},

		// empty paragraphs are saved as &nbsp; and come back empty
		{"empty paragraph", "uno\n\n&nbsp;\n\ndos", []Para{
			{Text: "uno"},
			{},
			{Text: "dos"},
		}},
		{"two empty paragraphs", "uno\n\n&nbsp;\n\n&nbsp;\n\ndos", []Para{
			{Text: "uno"},
			{},
			{},
			{Text: "dos"},
		}},
		{"empty paragraph inside a fence", "::: center\n\n&nbsp;\n\n:::", []Para{
			{Align: AlignCenter},
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
