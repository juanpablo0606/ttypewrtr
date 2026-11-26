package doc

import "testing"

func TestStyleValues(t *testing.T) {
	tests := []struct {
		style Style
		want  Style
	}{
		{Plain, 0},
		{Heading, 1},
		{Bold, 2},
		{Italic, 4},
		{Quote, 8},
		{Code, 16},
		{Marker, 32},
	}
	for _, tt := range tests {
		if tt.style != tt.want {
			t.Errorf("style has value %d, want %d", tt.style, tt.want)
		}
	}
}

func TestStyleString(t *testing.T) {
	tests := []struct {
		style Style
		want  string
	}{
		{Plain, "Plain"},
		{Heading, "Heading"},
		{Bold, "Bold"},
		{Italic, "Italic"},
		{Quote, "Quote"},
		{Code, "Code"},
		{Marker, "Marker"},
		{Heading | Bold, "Heading|Bold"},
		{Bold | Italic, "Bold|Italic"},
		{Italic | Bold, "Bold|Italic"}, // always in the order of the constants
		{Heading | Bold | Italic | Quote, "Heading|Bold|Italic|Quote"},
	}
	for _, tt := range tests {
		if got := tt.style.String(); got != tt.want {
			t.Errorf("Style(%d).String() = %q, want %q", uint8(tt.style), got, tt.want)
		}
	}
}

func TestStyleHas(t *testing.T) {
	s := Heading | Bold

	if !s.Has(Heading) {
		t.Error("(Heading|Bold).Has(Heading) = false, want true")
	}
	if !s.Has(Bold) {
		t.Error("(Heading|Bold).Has(Bold) = false, want true")
	}
	if s.Has(Italic) {
		t.Error("(Heading|Bold).Has(Italic) = true, want false")
	}
}
