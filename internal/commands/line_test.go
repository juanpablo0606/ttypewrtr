package commands

import (
	"testing"
)

func TestToggleHeading(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		level int
		want  string
	}{
		{"add h1", "hola", 1, "# hola"},
		{"remove", "# hola", 1, "hola"},
		{"Replace", "## hola", 1, "# hola"},
		{"Replace", "# hola", 3, "### hola"},
		{"empty line", "", 2, "## "},
		{"no space", "#hashtag", 1, "# #hashtag"},
		{"empty h1", "#", 1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToggleHeading(tt.line, tt.level)
			if got != tt.want {
				t.Errorf("ToggleHeading(%q, %d) = %q, want %q", tt.line, tt.level, got, tt.want)
			}
		})
	}
}

/*
func TestTogglePrefix(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		prefix stringish
		want  string
	}{
		{"add", "hola", , "# hola"},
		{"remove", "# hola", 1, "hola"},
		{"Replace", "## hola", 1, "# hola"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToggleHeading(tt.line, tt.level)
			if got != tt.want {
				t.Errorf("ToggleHeading(%q, %d) = %q, want %q", tt.line, tt.level, got, tt.want)
			}
		})
	}
}*/
