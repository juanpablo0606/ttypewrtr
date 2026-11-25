package config

import (
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeConfig writes content to a config.toml in a temporary folder
// and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

// --- Default ---

func TestDefaultHasSaveAndQuit(t *testing.T) {
	keys := Default().Keys
	for _, name := range []string{"save", "quit"} {
		if keys[name] == "" {
			t.Errorf("Default() has no key for %q", name)
		}
	}
}

// Default must build a new map on every call. If two configs shared the
// same map, changing one would silently change the other.
func TestDefaultReturnsIndependentMaps(t *testing.T) {
	a := Default()
	b := Default()

	a.Keys["bold"] = "changed"

	if b.Keys["bold"] == "changed" {
		t.Error("changing one Default() changed another: they share the same map")
	}
}

// --- Path ---

func TestPathEndsInTtypewrtrFolder(t *testing.T) {
	path, err := Path()
	if err != nil {
		t.Skipf("no user config dir in this environment: %v", err)
	}

	want := filepath.Join("ttypewrtr", "config.toml")
	if !strings.HasSuffix(path, want) {
		t.Errorf("Path() = %q, want it to end in %q", path, want)
	}
}

// --- Load ---

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.toml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load on a missing file returned error: %v", err)
	}
	if !maps.Equal(cfg.Keys, Default().Keys) {
		t.Errorf("Load on a missing file = %v, want the defaults %v", cfg.Keys, Default().Keys)
	}
}

func TestLoadEmptyFileReturnsDefaults(t *testing.T) {
	path := writeConfig(t, "")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load on an empty file returned error: %v", err)
	}
	if !maps.Equal(cfg.Keys, Default().Keys) {
		t.Errorf("Load on an empty file = %v, want the defaults", cfg.Keys)
	}
}

// Changing one shortcut must keep all the others at their default value.
func TestLoadMergesWithDefaults(t *testing.T) {
	path := writeConfig(t, `
[keys]
bold = "ctrl+b"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if got := cfg.Keys["bold"]; got != "ctrl+b" {
		t.Errorf("bold = %q, want %q", got, "ctrl+b")
	}

	want := Default().Keys
	want["bold"] = "ctrl+b"
	if !maps.Equal(cfg.Keys, want) {
		t.Errorf("Load = %v, want %v", cfg.Keys, want)
	}
}

func TestLoadComments(t *testing.T) {
	path := writeConfig(t, `
# my shortcuts
[keys]
quit = "alt+x" # VS Code eats ctrl+q
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got := cfg.Keys["quit"]; got != "alt+x" {
		t.Errorf("quit = %q, want %q", got, "alt+x")
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"broken TOML", "[keys\nbold = "},
		{"misspelled section", "[kyes]\nbold = \"alt+b\""},
		{"unknown top-level setting", "theme = \"dark\""},
		{"wrong value type", "[keys]\nbold = 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeConfig(t, tt.content)

			_, err := Load(path)
			if err == nil {
				t.Fatal("Load should return an error")
			}
			// The message must tell the user which file to fix.
			if !strings.Contains(err.Error(), path) {
				t.Errorf("error %q does not mention the file %q", err, path)
			}
		})
	}
}
