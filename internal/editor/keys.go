package editor

import "github.com/juanpablo0606/ttypewrtr/internal/commands"

type action func(m *Model)

var keymap = map[string]action{
	"alt+1": func(m *Model) { m.applyLine(heading(1)) },
	"alt+2": func(m *Model) { m.applyLine(heading(2)) },
	"alt+3": func(m *Model) { m.applyLine(heading(3)) },
	"alt+4": func(m *Model) { m.applyLine(heading(4)) },
	"alt+5": func(m *Model) { m.applyLine(heading(5)) },
	"alt+6": func(m *Model) { m.applyLine(heading(6)) },
	"alt+b": func(m *Model) { m.applyWord("**") },
	"alt+i": func(m *Model) { m.applyWord("*") },
	"alt+q": func(m *Model) { m.applyLine(prefix(">")) },
	"alt+l": func(m *Model) { m.applyLine(prefix("-")) },
}

func heading(level int) func(string) string {
	return func(s string) string {
		return commands.ToggleHeading(s, level)
	}
}

func prefix(prefix string) func(string) string {
	return func(s string) string {
		return commands.TogglePrefix(s, prefix)
	}
}
