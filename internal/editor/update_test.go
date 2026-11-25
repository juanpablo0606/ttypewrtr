package editor

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// Keys used in these tests, built the same way Bubble Tea builds them
// when you press them in a real terminal.
var (
	keyAltB  = tea.KeyPressMsg{Code: 'b', Mod: tea.ModAlt}
	keyAlt1  = tea.KeyPressMsg{Code: '1', Mod: tea.ModAlt}
	keyAltN  = tea.KeyPressMsg{Code: 'n', Mod: tea.ModAlt}
	keyAltX  = tea.KeyPressMsg{Code: 'x', Mod: tea.ModAlt}
	keyCtrlS = tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl}
	keyCtrlQ = tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl}
	keyCtrlC = tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	keyA     = tea.KeyPressMsg{Code: 'a', Text: "a"}
)

// press sends one key to Update and returns the new model and command.
func press(t *testing.T, m Model, key tea.KeyPressMsg) (Model, tea.Cmd) {
	t.Helper()

	next, cmd := m.Update(key)
	nm, ok := next.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want editor.Model", next)
	}
	return nm, cmd
}

// isQuit reports whether cmd is the quit command.
func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// Make sure the test keys produce the strings used in the config.
// If this fails, the other tests in this file are not testing what
// they claim to.
func TestKeyStrings(t *testing.T) {
	tests := []struct {
		key  tea.KeyPressMsg
		want string
	}{
		{keyAltB, "alt+b"},
		{keyAlt1, "alt+1"},
		{keyCtrlS, "ctrl+s"},
		{keyCtrlQ, "ctrl+q"},
		{keyCtrlC, emergencyQuitKey},
		{keyA, "a"},
	}
	for _, tt := range tests {
		if got := tt.key.String(); got != tt.want {
			t.Errorf("key.String() = %q, want %q", got, tt.want)
		}
	}
}

// --- formatting shortcuts ---

func TestUpdateFormatShortcuts(t *testing.T) {
	m := newTestModel(t, "hola", 0, 2)

	m, _ = press(t, m, keyAltB)
	if got := m.ta.Value(); got != "**hola**" {
		t.Errorf("after alt+b: %q, want %q", got, "**hola**")
	}

	m, _ = press(t, m, keyAlt1)
	if got := m.ta.Value(); got != "# **hola**" {
		t.Errorf("after alt+1: %q, want %q", got, "# **hola**")
	}
}

func TestUpdateUsesConfiguredKeys(t *testing.T) {
	m := newTestModelWithKeys(t, "hola", 0, 2, map[string]string{
		"bold": "alt+n",
		"save": "ctrl+s",
		"quit": "ctrl+q",
	})

	m, _ = press(t, m, keyAltB)
	if got := m.ta.Value(); got != "hola" {
		t.Errorf("alt+b is not configured, but changed the text to %q", got)
	}

	m, _ = press(t, m, keyAltN)
	if got := m.ta.Value(); got != "**hola**" {
		t.Errorf("after alt+n (configured bold): %q, want %q", got, "**hola**")
	}
}

// Keys without a shortcut must reach the textarea as normal typing.
func TestUpdateTypingReachesTextarea(t *testing.T) {
	m := newTestModel(t, "hol", 0, 3)

	m, _ = press(t, m, keyA)
	if got := m.ta.Value(); got != "hola" {
		t.Errorf("after typing a: %q, want %q", got, "hola")
	}
}

// --- save ---

func TestUpdateSave(t *testing.T) {
	m := newTestModel(t, "hola", 0, 4)
	m, _ = press(t, m, keyA)

	m, cmd := press(t, m, keyCtrlS)
	if cmd != nil {
		t.Error("save should not return a command")
	}
	if m.status != "Saved!" {
		t.Errorf("status = %q, want %q", m.status, "Saved!")
	}
	if m.hasChanges() {
		t.Error("after saving, hasChanges() should be false")
	}
}

// The status message must disappear on the next key.
func TestUpdateStatusClearsOnNextKey(t *testing.T) {
	m := newTestModel(t, "hola", 0, 4)
	m, _ = press(t, m, keyCtrlS)

	m, _ = press(t, m, keyA)
	if m.status != "" {
		t.Errorf("status after typing = %q, want empty", m.status)
	}
}

// --- quit ---

func TestUpdateQuitWithoutChanges(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{keyCtrlQ, keyCtrlC} {
		m := newTestModel(t, "hola", 0, 0)

		_, cmd := press(t, m, key)
		if !isQuit(cmd) {
			t.Errorf("%s without changes should quit at once", key.String())
		}
	}
}

func TestUpdateQuitAsksFirstWithChanges(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{keyCtrlQ, keyCtrlC} {
		t.Run(key.String(), func(t *testing.T) {
			m := newTestModel(t, "hola", 0, 4)
			m, _ = press(t, m, keyA)

			m, cmd := press(t, m, key)
			if isQuit(cmd) {
				t.Fatal("first quit with changes should warn, not quit")
			}
			if m.status == "" {
				t.Error("first quit with changes should show a warning")
			}

			_, cmd = press(t, m, key)
			if !isQuit(cmd) {
				t.Error("second quit should quit")
			}
		})
	}
}

// Any change after the warning must cancel it, or the user would lose
// that change without being warned again.
func TestUpdateQuitWarningResetsAfterEdit(t *testing.T) {
	m := newTestModel(t, "hola mundo", 0, 2)
	m, _ = press(t, m, keyAltB) // change 1

	m, _ = press(t, m, keyCtrlQ) // warning
	m, _ = press(t, m, keyAlt1)  // change 2, after the warning

	_, cmd := press(t, m, keyCtrlQ)
	if isQuit(cmd) {
		t.Error("quit after a new change should warn again, not quit")
	}
}

// ctrl+c quits even when the config moves quit to another key.
func TestUpdateEmergencyQuitAlwaysWorks(t *testing.T) {
	m := newTestModelWithKeys(t, "hola", 0, 0, map[string]string{
		"save": "ctrl+s",
		"quit": "alt+x",
	})

	_, cmd := press(t, m, keyCtrlC)
	if !isQuit(cmd) {
		t.Error("ctrl+c should quit even if quit is configured to alt+x")
	}

	_, cmd = press(t, m, keyAltX)
	if !isQuit(cmd) {
		t.Error("alt+x is configured as quit and should quit")
	}
}
