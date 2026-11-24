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
}

func New(path string) (Model, error) {

	m := Model{}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		//si está vacio no necesitamos cargar nada.
	} else if err != nil {
		//imprimir error al abrir el archivo
		return m, err
	}

	m.path = path

	content := string(data)

	m.ta = textarea.New()
	m.ta.SetValue(content)
	m.savedContent = content
	m.ta.Focus()

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+q":
			if m.hasChanges() && !m.quitWarned {
				m.status = "Unsaved changes. Ctrl+S to save, Ctrl+Q to quit"
				m.quitWarned = true
				return m, nil
			}
			return m, tea.Quit
		case "ctrl+s":
			err := m.save()
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			m.status = "Saved!"
			m.savedContent = m.ta.Value()
			return m, nil

		}
		m.quitWarned = false
		m.status = ""

		//key map
		if act, ok := keymap[msg.String()]; ok {
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
	saved := ""
	if m.hasChanges() {
		saved = "[+]"
	}
	statusBar := fmt.Sprintf("%s %s %s", m.path, saved, m.status)
	return tea.NewView(statusBar + "\n" + m.ta.View())
}

func (m Model) save() error {
	//0644 -> puedo leer y escribir, los demás solo leer
	return os.WriteFile(filepath.Base(m.path), []byte(m.ta.Value()), 0644)
}

func (m Model) hasChanges() bool {
	return m.savedContent != m.ta.Value()
}
