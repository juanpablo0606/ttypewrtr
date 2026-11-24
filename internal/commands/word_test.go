package commands

import (
	"testing"
)

func TestWrapWord(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		col     int
		marker  string
		want    string
		wantCol int
	}{
		{"cursor at middle", "hola mundo", 2, "**", "**hola** mundo", 4},
		{"cursor at start", "hola mundo", 0, "**", "**hola** mundo", 2},
		{"cursor after word", "hola mundo", 4, "**", "**hola** mundo", 6},
		{"cursor at second word", "hola mundo", 8, "*", "hola *mundo*", 9},
		{"special char", "cañón", 2, "**", "**cañón**", 4},
		{"cursor out of word", "hola ", 5, "**", "hola ****", 7},
		{"empty line", "", 0, "*", "**", 1},
		{"cursor far", "hola", 99, "**", "**hola**", 6},
	}

	for _, tt := range tests {
		got, gotCol := WrapWord(tt.line, tt.col, tt.marker)
		if got != tt.want || gotCol != tt.wantCol {
			t.Errorf("WrapWord(%q, %d, %q) = (%q, %d), want (%q, %d)",
				tt.line, tt.col, tt.marker, got, gotCol, tt.want, tt.wantCol)
		}
	}
}
