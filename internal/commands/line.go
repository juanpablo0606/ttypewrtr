package commands

import (
	"strings"
)

// transforms into heading level 1-6
// if same level > remove
// if other level > replace
func ToggleHeading(line string, level int) string {
	hLevel := headingLevel(line)
	if hLevel == 0 {
		return addHeading(line, level)
	}
	trim := strings.TrimLeft(line, "#")
	trim = strings.TrimLeft(trim, " ")
	if hLevel == level {
		return trim
	}
	return addHeading(trim, level)

}

// add prefix to line or removes it if already there
// Quotes ("> ") and lists("- ")

func TogglePrefix(line, prefix string) string {
	if strings.HasPrefix(line, prefix) {
		return strings.TrimPrefix(line, prefix)
	}
	return prefix + line
}

// check how many "#" in a row at the start
func headingLevel(line string) int {
	rest := strings.TrimLeft(line, "#")
	level := len(line) - len(rest)

	//check if it's a heading (heading needs a space after "#")
	if strings.HasPrefix(rest, " ") && level <= 6 {
		return level
	}
	return 0
}

func addHeading(line string, level int) string {
	heading := strings.Repeat("#", level)
	heading += " "
	return heading + line
}
