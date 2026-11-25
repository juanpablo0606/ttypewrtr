package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/juanpablo0606/ttypewrtr/internal/config"
)

// newTestModel creates a Model with the given content and puts the cursor
// at (row, col). The file lives in a temporary folder that Go deletes
// automatically when the test ends.
func newTestModel(t *testing.T, content string, row, col int) Model {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := New(path, config.Default().Keys)
	if err != nil {
		t.Fatalf("New(%q) returned error: %v", path, err)
	}
	m.ta.SetWidth(80)
	m.ta.SetHeight(20)
	m.setCursor(row, col)
	return m
}

// assertState checks the full text and the cursor position.
func assertState(t *testing.T, m Model, wantText string, wantRow, wantCol int) {
	t.Helper()

	if got := m.ta.Value(); got != wantText {
		t.Errorf("text = %q, want %q", got, wantText)
	}
	if m.ta.Line() != wantRow || m.ta.Column() != wantCol {
		t.Errorf("cursor = (%d, %d), want (%d, %d)", m.ta.Line(), m.ta.Column(), wantRow, wantCol)
	}
}

// --- New ---

func TestNewExistingFile(t *testing.T) {
	m := newTestModel(t, "hola\nmundo", 0, 0)

	if got := m.ta.Value(); got != "hola\nmundo" {
		t.Errorf("content = %q, want %q", got, "hola\nmundo")
	}
	if m.hasChanges() {
		t.Error("a freshly opened file must not have changes")
	}
}

func TestNewMissingFileStartsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.md")

	m, err := New(path, config.Default().Keys)
	if err != nil {
		t.Fatalf("New on a missing file returned error: %v", err)
	}
	if m.ta.Value() != "" {
		t.Errorf("content = %q, want empty", m.ta.Value())
	}
	if m.hasChanges() {
		t.Error("a new empty file must not have changes")
	}
}

func TestNewUnreadablePathReturnsError(t *testing.T) {
	// A directory cannot be read as a file.
	_, err := New(t.TempDir(), config.Default().Keys)
	if err == nil {
		t.Error("New on a directory should return an error")
	}
}

// --- save ---

// Regression test: save must write to the full path, not to the
// current working directory.
func TestSaveWritesToFullPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "notes.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}

	m, err := New(path, config.Default().Keys)
	if err != nil {
		t.Fatal(err)
	}
	m.ta.SetValue("contenido")

	if err := m.save(); err != nil {
		t.Fatalf("save() returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not found at %q: %v", path, err)
	}
	if string(data) != "contenido" {
		t.Errorf("file content = %q, want %q", data, "contenido")
	}
}

// --- hasChanges ---

func TestHasChanges(t *testing.T) {
	m := newTestModel(t, "hola", 0, 4)

	m.ta.InsertString("!")
	if !m.hasChanges() {
		t.Error("after typing, hasChanges() should be true")
	}

	m.ta.SetValue("hola")
	if m.hasChanges() {
		t.Error("after restoring the original text, hasChanges() should be false")
	}
}

// --- applyLine (requirements R1–R5, R7–R9 of Part D) ---

func TestApplyLine(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		row, col int
		fn       func(string) string
		wantText string
		wantRow  int
		wantCol  int
	}{
		{
			name:    "R1/R3: only the cursor line changes, cursor stays on it",
			content: "uno\ndos\ntres", row: 1, col: 1,
			fn:       heading(1),
			wantText: "uno\n# dos\ntres", wantRow: 1, wantCol: 3,
		},
		{
			name:    "R2: trailing newline is kept",
			content: "uno\ndos\n", row: 0, col: 0,
			fn:       heading(1),
			wantText: "# uno\ndos\n", wantRow: 0, wantCol: 2,
		},
		{
			name:    "R4: cursor stays on the same letter when adding",
			content: "hola", row: 0, col: 2,
			fn:       heading(2),
			wantText: "## hola", wantRow: 0, wantCol: 5,
		},
		{
			name:    "R4: cursor stays on the same letter when removing",
			content: "# hola", row: 0, col: 4,
			fn:       heading(1),
			wantText: "hola", wantRow: 0, wantCol: 2,
		},
		{
			name:    "R5: cursor on the removed prefix goes to column 0",
			content: "# hola", row: 0, col: 1,
			fn:       heading(1),
			wantText: "hola", wantRow: 0, wantCol: 0,
		},
		{
			name:    "replace heading level",
			content: "# hola", row: 0, col: 4,
			fn:       heading(3),
			wantText: "### hola", wantRow: 0, wantCol: 6,
		},
		{
			name:    "quote",
			content: "uno\ndos", row: 1, col: 0,
			fn:       prefix("> "),
			wantText: "uno\n> dos", wantRow: 1, wantCol: 2,
		},
		{
			name:    "list",
			content: "uno\ndos", row: 0, col: 3,
			fn:       prefix("- "),
			wantText: "- uno\ndos", wantRow: 0, wantCol: 5,
		},
		{
			name:    "R7: long wrapped line above the cursor",
			content: strings.Repeat("palabra ", 30) + "\ndos", row: 1, col: 1,
			fn:       heading(1),
			wantText: strings.Repeat("palabra ", 30) + "\n# dos", wantRow: 1, wantCol: 3,
		},
		{
			name:    "R8: empty document",
			content: "", row: 0, col: 0,
			fn:       heading(1),
			wantText: "# ", wantRow: 0, wantCol: 2,
		},
		{
			name:    "R8: last line",
			content: "uno\ndos\ntres", row: 2, col: 4,
			fn:       heading(1),
			wantText: "uno\ndos\n# tres", wantRow: 2, wantCol: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel(t, tt.content, tt.row, tt.col)
			m.applyLine(tt.fn)
			assertState(t, m, tt.wantText, tt.wantRow, tt.wantCol)
		})
	}
}

// R9: toggling twice restores the text, so the [+] marker goes away.
func TestApplyLineTwiceHasNoChanges(t *testing.T) {
	m := newTestModel(t, "uno\ndos", 1, 1)

	m.applyLine(heading(1))
	if !m.hasChanges() {
		t.Error("after one toggle, hasChanges() should be true")
	}

	m.applyLine(heading(1))
	if m.hasChanges() {
		t.Error("after toggling twice, hasChanges() should be false")
	}
	assertState(t, m, "uno\ndos", 1, 1)
}

// --- applyWord (requirement R6) ---

func TestApplyWord(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		row, col int
		marker   string
		wantText string
		wantRow  int
		wantCol  int
	}{
		{
			name:    "bold on second line",
			content: "uno\nhola mundo", row: 1, col: 2,
			marker:   "**",
			wantText: "uno\n**hola** mundo", wantRow: 1, wantCol: 4,
		},
		{
			name:    "italic, other lines untouched",
			content: "uno\ndos\ntres", row: 1, col: 1,
			marker:   "*",
			wantText: "uno\n*dos*\ntres", wantRow: 1, wantCol: 2,
		},
		{
			name:    "no word: empty markers with cursor in the middle",
			content: "hola ", row: 0, col: 5,
			marker:   "**",
			wantText: "hola ****", wantRow: 0, wantCol: 7,
		},
		{
			name:    "empty document",
			content: "", row: 0, col: 0,
			marker:   "*",
			wantText: "**", wantRow: 0, wantCol: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel(t, tt.content, tt.row, tt.col)
			m.applyWord(tt.marker)
			assertState(t, m, tt.wantText, tt.wantRow, tt.wantCol)
		})
	}
}

// --- keymap ---

// Every shortcut must point to an action. Catches typos like a nil entry.
func TestKeymapHasNoNilActions(t *testing.T) {
	for key, act := range actions {
		if act == nil {
			t.Errorf("keymap[%q] is nil", key)
		}
	}
}

// Shortcuts that would make it impossible to save or quit must not be
// used for formatting.
func TestKeymapDoesNotShadowReservedKeys(t *testing.T) {
	reserved := []string{"ctrl+s", "ctrl+q", "ctrl+c"}
	for _, key := range reserved {
		if _, ok := actions[key]; ok {
			t.Errorf("keymap uses reserved key %q", key)
		}
	}
}
