package commands

import "testing"

func TestHeadingLevel(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int
	}{
		{"plain text", "hola", 0},
		{"h1", "# hola", 1},
		{"h3", "### hola", 3},
		{"h6", "###### hola", 6},
		{"seven hashes is not a heading", "####### hola", 0},
		{"no space is not a heading", "#hola", 0},
		{"lonely hash is an empty h1", "#", 1},
		{"only hashes is an empty heading", "###", 3},
		{"empty line", "", 0},
		{"hash in the middle", "hola # mundo", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := headingLevel(tt.line)
			if got != tt.want {
				t.Errorf("headingLevel(%q) = %d, want %d", tt.line, got, tt.want)
			}
		})
	}
}

func TestToggleHeading(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		level int
		want  string
	}{
		// add
		{"add h1", "hola", 1, "# hola"},
		{"add h3", "hola", 3, "### hola"},
		{"add to empty line", "", 2, "## "},
		{"no space is not a heading, so add", "#hashtag", 1, "# #hashtag"},
		{"seven hashes is not a heading, so add", "####### hola", 1, "# ####### hola"},

		// remove (same level)
		{"remove h1", "# hola", 1, "hola"},
		{"remove h2", "## hola", 2, "hola"},
		{"remove empty h1", "#", 1, ""},
		{"remove h1 with only a space", "# ", 1, ""},
		{"remove h1 with extra spaces", "#   hola", 1, "hola"},

		// replace (different level)
		{"replace h2 with h1", "## hola", 1, "# hola"},
		{"replace h1 with h3", "# hola", 3, "### hola"},
		{"replace h6 with h1", "###### hola", 1, "# hola"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToggleHeading(tt.line, tt.level)
			if got != tt.want {
				t.Errorf("ToggleHeading(%q, %d) = %q, want %q", tt.line, tt.level, got, tt.want)
			}
		})
	}
}

// Applying the same heading twice must give back the original line
// (as long as the original line was not already a heading).
func TestToggleHeadingTwiceIsIdentity(t *testing.T) {
	lines := []string{"hola", "", "hola mundo", "cañón", "#hashtag", "> cita"}

	for _, line := range lines {
		for level := 1; level <= 6; level++ {
			once := ToggleHeading(line, level)
			twice := ToggleHeading(once, level)
			if twice != line {
				t.Errorf("ToggleHeading twice(%q, %d) = %q, want %q (once = %q)", line, level, twice, line, once)
			}
		}
	}
}

func TestTogglePrefix(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		prefix string
		want   string
	}{
		// quotes
		{"add quote", "hola", "> ", "> hola"},
		{"remove quote", "> hola", "> ", "hola"},
		{"add quote to empty line", "", "> ", "> "},
		{"remove quote leaving empty line", "> ", "> ", ""},
		{"quote without space is not recognized", ">hola", "> ", "> >hola"},

		// lists
		{"add list", "hola", "- ", "- hola"},
		{"remove list", "- hola", "- ", "hola"},
		{"dash without space is not a list", "-hola", "- ", "- -hola"},

		// combining different prefixes
		{"quote a list item", "- hola", "> ", "> - hola"},
		{"list a quote", "> hola", "- ", "- > hola"},

		// known limitation: indentation is ignored
		{"indented list (known limitation)", "  - hola", "- ", "-   - hola"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TogglePrefix(tt.line, tt.prefix)
			if got != tt.want {
				t.Errorf("TogglePrefix(%q, %q) = %q, want %q", tt.line, tt.prefix, got, tt.want)
			}
		})
	}
}

// Applying the same prefix twice must always give back the original line.
func TestTogglePrefixTwiceIsIdentity(t *testing.T) {
	lines := []string{"hola", "", "> hola", "- hola", ">hola", "  - hola", "cañón"}
	prefixes := []string{"> ", "- "}

	for _, line := range lines {
		for _, p := range prefixes {
			twice := TogglePrefix(TogglePrefix(line, p), p)
			if twice != line {
				t.Errorf("TogglePrefix twice(%q, %q) = %q, want %q", line, p, twice, line)
			}
		}
	}
}
