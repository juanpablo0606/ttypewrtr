package editor

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/juanpablo0606/ttypewrtr/internal/commands"
)

const emergencyQuitKey = "ctrl+c"

// action is something the user can trigger with a key.
type action func(m *Model) tea.Cmd

var actions = map[string]action{
	"save": func(m *Model) tea.Cmd { return m.saveFile() },
	"quit": func(m *Model) tea.Cmd { return m.quit() },

	"h1": lineAction(heading(1)),
	"h2": lineAction(heading(2)),
	"h3": lineAction(heading(3)),
	"h4": lineAction(heading(4)),
	"h5": lineAction(heading(5)),
	"h6": lineAction(heading(6)),

	"bold":   wordAction("**"),
	"italic": wordAction("*"),
	"quote":  lineAction(prefix("> ")),
	"list":   lineAction(prefix("- ")),
}

func heading(level int) func(string) string {
	return func(s string) string {
		return commands.ToggleHeading(s, level)
	}
}

func prefix(p string) func(string) string {
	return func(s string) string {
		return commands.TogglePrefix(s, p)
	}
}

func buildKeymap(keys map[string]string) (map[string]string, error) {
	keymap := make(map[string]string)
	for name, key := range keys {
		if _, ok := actions[name]; !ok {
			return nil, fmt.Errorf("unknown action %q", name)
		}
		if key == "" {
			return nil, fmt.Errorf("action %q has no key", name)
		}
		if key == emergencyQuitKey && name != "quit" {
			return nil, fmt.Errorf("key %q is reserved for quitting", key)
		}
		if other, used := keymap[key]; used {
			return nil, fmt.Errorf("key %q is used by both %q and %q", key, name, other)
		}
		keymap[key] = name
	}
	return keymap, nil
}

// lineAction returns an action that transforms the cursor line with fn.
func lineAction(fn func(string) string) action {
	return func(m *Model) tea.Cmd {
		m.applyLine(fn)
		return nil
	}
}

// wordAction returns an action that toggles marker around the cursor word.
func wordAction(marker string) action {
	return func(m *Model) tea.Cmd {
		m.applyWord(marker)
		return nil
	}
}

// saveFile writes the file and reports the result in the status bar.
func (m *Model) saveFile() tea.Cmd {
	err := m.save()
	if err != nil {
		m.status = err.Error()
		return nil
	}
	m.status = "Saved!"
	m.savedContent = m.ta.Value()
	return nil
}

// quit exits, but asks for confirmation first if there are unsaved changes.
func (m *Model) quit() tea.Cmd {
	if m.hasChanges() && !m.quitWarned {
		m.status = "Unsaved changes. Ctrl+S to save, Ctrl+Q to quit"
		m.quitWarned = true
		return nil
	}
	return tea.Quit
}
