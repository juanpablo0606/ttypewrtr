package commands

import "unicode"

// WrapWord toggles marker ("*" or "**") around the word under the cursor.
// If there is no word, it inserts an empty pair of markers.
// It returns the new line and the new cursor column.
func WrapWord(line string, col int, marker string) (string, int) {
	r := []rune(line)
	col = clamp(col, 0, len(r))
	start, end := wordBounds(r, col)
	m := []rune(marker)

	stars := surrounding(r, start, end, m[0])
	if isApplied(stars, len(m)) {
		return unwrap(r, start, end, len(m)), col - len(m)
	}
	return wrap(r, start, end, m), col + len(m)
}

func isWordChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}

// clamp limits v to the range [lo, hi].
func clamp(v, lo, hi int) int {
	if v < lo {
		v = lo
	} else if v > hi {
		v = hi
	}
	return v
}

// wordBounds returns where the word around col starts and ends.
// If col is not touching a word, start == end == col.
func wordBounds(r []rune, col int) (start, end int) {
	start = col
	for start > 0 && isWordChar(r[start-1]) {
		start--
	}
	end = col
	for end < len(r) && isWordChar(r[end]) {
		end++
	}
	return start, end
}

// surrounding counts how many ch runes surround r[start:end] on both sides.
// Only stars present on BOTH sides count, so it returns the smaller count.
func surrounding(r []rune, start, end int, ch rune) int {
	left := 0
	for start-left > 0 && r[start-left-1] == ch {
		left++
	}
	right := 0
	for end+right < len(r) && r[end+right] == ch {
		right++
	}
	return min(left, right)

}

// isApplied reports whether a marker of length n is already present,
// given the number of stars around the word.
//
//	italic (n=1): present with an odd count (1 or 3)
//	bold   (n=2): present with 2 or more
func isApplied(stars, n int) bool {
	italic := stars%2 == 1 // impar: 1, 3, 5...
	bold := stars >= 2     // 2 o más

	if n == 1 {
		return italic
	}
	return bold
}

// wrap returns r with marker inserted around r[start:end].
func wrap(r []rune, start, end int, marker []rune) string {
	return string(r[:start]) + string(marker) + string(r[start:end]) + string(marker) + string(r[end:])
}

// unwrap returns r with n runes removed on each side of r[start:end].
func unwrap(r []rune, start, end, n int) string {
	return string(r[:start-n]) + string(r[start:end]) + string(r[end+n:])
}
