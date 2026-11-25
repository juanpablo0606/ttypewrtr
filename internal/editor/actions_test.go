package editor

import (
	"testing"

	"github.com/juanpablo0606/ttypewrtr/internal/config"
)

// Every action in the default config must exist in actions, and every
// action must have a default key. If you add an action and forget one of
// the two, this test fails.
func TestDefaultConfigMatchesActions(t *testing.T) {
	keys := config.Default().Keys

	for name := range keys {
		if _, ok := actions[name]; !ok {
			t.Errorf("config.Default() has %q, but actions does not", name)
		}
	}
	for name := range actions {
		if _, ok := keys[name]; !ok {
			t.Errorf("actions has %q, but config.Default() gives it no key", name)
		}
	}
}

func TestActionsHaveNoNilEntries(t *testing.T) {
	for name, act := range actions {
		if act == nil {
			t.Errorf("actions[%q] is nil", name)
		}
	}
}

func TestBuildKeymapDefaults(t *testing.T) {
	keys := config.Default().Keys

	km, err := buildKeymap(keys)
	if err != nil {
		t.Fatalf("buildKeymap(defaults) returned error: %v", err)
	}

	// Every action must be reachable from its key.
	for name, key := range keys {
		if km[key] != name {
			t.Errorf("keymap[%q] = %q, want %q", key, km[key], name)
		}
	}
	if len(km) != len(keys) {
		t.Errorf("keymap has %d entries, want %d", len(km), len(keys))
	}
}

func TestBuildKeymapInverts(t *testing.T) {
	km, err := buildKeymap(map[string]string{
		"bold": "alt+n",
		"quit": "alt+x",
	})
	if err != nil {
		t.Fatalf("buildKeymap returned error: %v", err)
	}

	if km["alt+n"] != "bold" {
		t.Errorf(`keymap["alt+n"] = %q, want "bold"`, km["alt+n"])
	}
	if km["alt+x"] != "quit" {
		t.Errorf(`keymap["alt+x"] = %q, want "quit"`, km["alt+x"])
	}
}

func TestBuildKeymapErrors(t *testing.T) {
	tests := []struct {
		name string
		keys map[string]string
	}{
		{"unknown action", map[string]string{"bolt": "alt+b"}},
		{"empty key", map[string]string{"quit": ""}},
		{"same key twice", map[string]string{"bold": "alt+b", "italic": "alt+b"}},
		{"emergency key on another action", map[string]string{"bold": emergencyQuitKey}},
		{"emergency key on save", map[string]string{"save": emergencyQuitKey}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km, err := buildKeymap(tt.keys)
			if err == nil {
				t.Errorf("buildKeymap(%v) = %v, want an error", tt.keys, km)
			}
		})
	}
}

// quit is the only action allowed to use the emergency key.
func TestBuildKeymapAllowsEmergencyKeyOnQuit(t *testing.T) {
	if _, err := buildKeymap(map[string]string{"quit": emergencyQuitKey}); err != nil {
		t.Errorf("quit = %q should be allowed, got error: %v", emergencyQuitKey, err)
	}
}

// A bad config must stop New, so the program never starts with broken keys.
func TestNewFailsOnBadKeys(t *testing.T) {
	path := t.TempDir() + "/test.md"

	if _, err := New(path, map[string]string{"bolt": "alt+b"}); err == nil {
		t.Error("New with an unknown action should return an error")
	}
}
