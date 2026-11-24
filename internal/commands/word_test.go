package commands

import "testing"

// In the comments, | marks the cursor position.

func TestWrapWord(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		col     int
		marker  string
		want    string
		wantCol int
	}{
		// wrap a word
		{"cursor in the middle: ho|la", "hola mundo", 2, "**", "**hola** mundo", 4},
		{"cursor at word start: |hola", "hola mundo", 0, "**", "**hola** mundo", 2},
		{"cursor right after word: hola|", "hola mundo", 4, "**", "**hola** mundo", 6},
		{"second word: hola mun|do", "hola mundo", 8, "*", "hola *mundo*", 9},
		{"accents and ñ: ca|ñón", "cañón", 2, "**", "**cañón**", 4},
		{"digits are part of the word: ab|c123", "abc123 x", 2, "*", "*abc123* x", 3},
		{"punctuation is not: ho|la, mundo", "hola, mundo", 2, "**", "**hola**, mundo", 4},

		// no word under the cursor: insert empty markers
		{"after a space: hola |", "hola ", 5, "**", "hola ****", 7},
		{"empty line", "", 0, "*", "**", 1},
		{"between two spaces: a | b", "a  b", 2, "*", "a ** b", 3},

		// invalid cursor positions are clamped
		{"cursor too far right", "hola", 99, "**", "**hola**", 6},
		{"negative cursor", "hola", -5, "**", "**hola**", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotCol := WrapWord(tt.line, tt.col, tt.marker)
			if got != tt.want || gotCol != tt.wantCol {
				t.Errorf("WrapWord(%q, %d, %q) = (%q, %d), want (%q, %d)",
					tt.line, tt.col, tt.marker, got, gotCol, tt.want, tt.wantCol)
			}
		})
	}
}

// TestWrapWordToggle describes the toggle behavior: if the word is already
// wrapped with the marker, WrapWord removes it.
//
// Bold (**) and italic (*) use the same character, so markdown decides by
// counting the stars around the word:
//
//	0 stars: plain       hola
//	1 star:  italic      *hola*
//	2 stars: bold        **hola**
//	3 stars: both        ***hola***
//
// Toggling italic adds or removes 1 star, toggling bold adds or removes 2.
func TestWrapWordToggle(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		col     int
		marker  string
		want    string
		wantCol int
	}{
		// remove the same marker
		{"remove bold: **ho|la**", "**hola** mundo", 4, "**", "hola mundo", 2},
		{"remove italic: *ho|la*", "*hola*", 3, "*", "hola", 2},
		{"remove bold, cursor right after word: **hola|**", "**hola**", 6, "**", "hola", 4},
		{"remove bold, accents: **ca|ñón**", "**cañón**", 4, "**", "cañón", 2},

		// stack bold and italic
		{"italic on bold: **ho|la**", "**hola**", 4, "*", "***hola***", 5},
		{"bold on italic: *ho|la*", "*hola*", 3, "**", "***hola***", 5},

		// remove one of the two from bold+italic
		{"remove italic from both: ***ho|la***", "***hola***", 5, "*", "**hola**", 4},
		{"remove bold from both: ***ho|la***", "***hola***", 5, "**", "*hola*", 3},

		// empty markers inserted by mistake are removed again
		{"remove empty bold: hola **|**", "hola ****", 7, "**", "hola ", 5},
		{"remove empty italic: *|*", "**", 1, "*", "", 0},

		// only the markers touching THIS word count
		{"other word is bold: **hola** mun|do", "**hola** mundo", 11, "**", "**hola** **mundo**", 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotCol := WrapWord(tt.line, tt.col, tt.marker)
			if got != tt.want || gotCol != tt.wantCol {
				t.Errorf("WrapWord(%q, %d, %q) = (%q, %d), want (%q, %d)",
					tt.line, tt.col, tt.marker, got, gotCol, tt.want, tt.wantCol)
			}
		})
	}
}

// Applying the same marker twice must give back the original line and cursor.
func TestWrapWordTwiceIsIdentity(t *testing.T) {
	tests := []struct {
		line string
		col  int
	}{
		{"hola mundo", 2},
		{"hola mundo", 8},
		{"cañón", 3},
		{"**hola**", 4},
		{"*hola*", 3},
		{"hola ", 5},
		{"", 0},
	}

	for _, tt := range tests {
		for _, marker := range []string{"*", "**"} {
			once, onceCol := WrapWord(tt.line, tt.col, marker)
			twice, twiceCol := WrapWord(once, onceCol, marker)
			if twice != tt.line || twiceCol != tt.col {
				t.Errorf("WrapWord twice(%q, %d, %q) = (%q, %d), want (%q, %d) (once = %q, %d)",
					tt.line, tt.col, marker, twice, twiceCol, tt.line, tt.col, once, onceCol)
			}
		}
	}
}
