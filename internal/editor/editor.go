package editor

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	ta           textarea.Model
	path         string
	savedContent string
	status       string
	quitWarned   bool //hay cambios sin guardar?
}

func New(path string) (Model, error) {

	m := Model{}
	m.path = path

	data, err := os.ReadFile(m.path)
	if errors.Is(err, fs.ErrNotExist) {
		//si está vacio no necesitamos cargar nada.
	} else if err != nil {
		//imprimir error al abrir el archivo
		return m, err
	}

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
		m.status = "[+]"
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			if !m.quitWarned {
				m.status = "¿Save before you leave?"
				m.quitWarned = true
				return m, nil
			}
			return m, tea.Quit
		case "ctrl+s":
			err := m.save()
			if err != nil {
				m.status = err.Error()
			}
			m.status = "Saved!"
			m.savedContent = m.ta.Value()
			return m, nil
		}
		m.quitWarned = false
	case tea.WindowSizeMsg:
		m.ta.SetWidth(msg.Width)
		m.ta.SetHeight(msg.Height - 5)
		return m, nil
	}

	var cmd tea.Cmd
	m.ta, cmd = m.ta.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {

	statusBar := fmt.Sprintf("%s %s", m.path, m.status)
	return tea.NewView(statusBar + "\n" + m.ta.View())
}

func (m Model) save() error {
	//0644 -> puedo leer y escribir, los demás solo leer
	return os.WriteFile(m.path, []byte(m.ta.Value()), 0644)
}
