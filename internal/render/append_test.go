package render

import (
	"slices"
	"testing"
)

func TestAppendSpan(t *testing.T) {
	tests := []struct {
		name  string
		spans []Span
		add   Span
		want  []Span
	}{
		{
			name:  "first span",
			spans: nil,
			add:   Span{0, 4, Plain},
			want:  []Span{{0, 4, Plain}},
		},
		{
			name:  "different style: append",
			spans: []Span{{0, 2, Marker}},
			add:   Span{2, 6, Heading},
			want:  []Span{{0, 2, Marker}, {2, 6, Heading}},
		},
		{
			name:  "same style right after: merge",
			spans: []Span{{2, 4, Bold}},
			add:   Span{4, 5, Bold},
			want:  []Span{{2, 5, Bold}},
		},
		{
			name:  "merge only touches the last span",
			spans: []Span{{0, 2, Marker}, {2, 4, Bold}},
			add:   Span{4, 5, Bold},
			want:  []Span{{0, 2, Marker}, {2, 5, Bold}},
		},
		{
			name:  "empty span: ignore",
			spans: []Span{{0, 2, Marker}},
			add:   Span{2, 2, Heading},
			want:  []Span{{0, 2, Marker}},
		},
		{
			name:  "empty span on empty list: ignore",
			spans: nil,
			add:   Span{0, 0, Marker},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendSpan(tt.spans, tt.add)
			if !slices.Equal(got, tt.want) {
				t.Errorf("appendSpan(%v, %v)\n got  %v\n want %v", tt.spans, tt.add, got, tt.want)
			}
		})
	}
}
