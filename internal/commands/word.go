package commands

import "unicode"

// col = position of cursor
func WrapWord(line string, col int, marker string) (string, int) {
	r := []rune(line)
	if col < 0 {
		col = 0
	} else if col > len(r) {
		col = len(r)
	}

	start := col
	for start > 0 && isWordChar(r[start-1]) {
		start--
	}
	end := col
	for end < len(r) && isWordChar(r[end]) {
		end++
	}

	markerLen := len([]rune(marker))
	newCol := col + markerLen

	//cursor is not touching any word
	if start == end {
		return string(r[:col]) + marker + string(r[col:]) + marker, newCol
	}
	return string(r[:start]) + marker + string(r[start:end]) + marker + string(r[end:]), newCol

}

func isWordChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}
