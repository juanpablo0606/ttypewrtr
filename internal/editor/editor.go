package editor

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

const statusBarHeight = 3

type Model struct {
	ta           textarea.Model
	path         string
	savedContent string
	status       string
	quitWarned   bool //hay cambios sin guardar?
	keymap       map[string]string
}

func New(path string, keys map[string]string) (Model, error) {
	km, err := buildKeymap(keys)
	if err != nil {
		return Model{}, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		//empty
	} else if err != nil {
		//err when opening
		return Model{}, err
	}

	content := string(data)

	ta := textarea.New()
	ta.SetValue(content)
	ta.Focus()

	return Model{
		ta:           ta,
		path:         path,
		savedContent: content,
		keymap:       km,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		name, bound := m.keymap[key]

		if key == emergencyQuitKey {
			name, bound = "quit", true
		}

		if name != "quit" {
			m.quitWarned = false
		}
		m.status = ""

		if bound {
			return m, actions[name](&m)
		}

		//key map
		if act, ok := actions[msg.String()]; ok {
			act(&m)
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.ta.SetWidth(msg.Width)
		m.ta.SetHeight(msg.Height - statusBarHeight)
		return m, nil
	}

	var cmd tea.Cmd
	m.ta, cmd = m.ta.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {
	modified := ""
	if m.hasChanges() {
		modified = "[+]"
	}
	statusBar := fmt.Sprintf("%s %s %s", filepath.Base(m.path), modified, m.status)
	return tea.NewView(statusBar + "\n" + m.ta.View())
}

func (m Model) save() error {
	//0644 -> puedo leer y escribir, los demás solo leer
	return os.WriteFile(m.path, []byte(m.ta.Value()), 0644)
}

func (m Model) hasChanges() bool {
	return m.savedContent != m.ta.Value()
}
