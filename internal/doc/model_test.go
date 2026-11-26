package doc

import "testing"

// The zero value of each type must mean "nothing special": a normal
// paragraph, with the document's default alignment.
func TestZeroValues(t *testing.T) {
	var p Para

	if p.Kind != KindBody {
		t.Errorf("zero Para has Kind %d, want KindBody (%d)", p.Kind, KindBody)
	}
	if p.Align != AlignDefault {
		t.Errorf("zero Para has Align %d, want AlignDefault (%d)", p.Align, AlignDefault)
	}
	if p.Level != 0 || p.Text != "" || p.Spans != nil {
		t.Errorf("zero Para = %+v, want everything empty", p)
	}
}

func TestKindValues(t *testing.T) {
	tests := []struct {
		kind Kind
		want int
	}{
		{KindBody, 0},
		{KindHeading, 1},
		{KindQuote, 2},
		{KindListItem, 3},
	}
	for _, tt := range tests {
		if int(tt.kind) != tt.want {
			t.Errorf("kind has value %d, want %d", tt.kind, tt.want)
		}
	}
}

func TestAlignValues(t *testing.T) {
	tests := []struct {
		align Align
		want  int
	}{
		{AlignDefault, 0},
		{AlignLeft, 1},
		{AlignCenter, 2},
		{AlignRight, 3},
		{AlignJustify, 4},
	}
	for _, tt := range tests {
		if int(tt.align) != tt.want {
			t.Errorf("align has value %d, want %d", tt.align, tt.want)
		}
	}
}

func TestDocText(t *testing.T) {
	tests := []struct {
		name string
		doc  Doc
		want string
	}{
		{"empty document", Doc{}, ""},
		{"one paragraph", Doc{Paras: []Para{{Text: "hola"}}}, "hola"},
		{
			name: "several paragraphs, one per line",
			doc:  Doc{Paras: []Para{{Text: "uno"}, {Text: "dos"}, {Text: "tres"}}},
			want: "uno\ndos\ntres",
		},
		{
			name: "empty paragraphs are empty lines",
			doc:  Doc{Paras: []Para{{Text: "uno"}, {}, {Text: "dos"}}},
			want: "uno\n\ndos",
		},
		{
			name: "kind, level and spans do not change the text",
			doc: Doc{Paras: []Para{
				{Kind: KindHeading, Level: 1, Text: "Capítulo"},
				{Text: "hola mundo", Spans: []Span{{5, 10, Bold}}},
			}},
			want: "Capítulo\nhola mundo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.doc.Text(); got != tt.want {
				t.Errorf("Text() = %q, want %q", got, tt.want)
			}
		})
	}
}
