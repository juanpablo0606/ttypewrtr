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
	ta    textarea.Model
	path  string
	saved bool
	err   error
}

func New(m Model, path string) (Model, error) {

	m.path = path

	data, err := os.ReadFile(m.path)
	if errors.Is(err, fs.ErrNotExist) {
		//don't know what to do.
	} else if err != nil {
		//imprimir error al abrir el archivo
		fmt.Println(err)
		os.Exit(1)
	}

	content := string(data)

	m.ta = textarea.New()
	m.ta.SetValue(content)
	m.ta.Focus()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return Model{}, err
	}
	return Model{}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			return m, tea.Quit
		case "ctrl+s":
			//me falta chequear error, a donde va el error?
			m.save()
			m.saved = true
			return m, nil
		}
	case tea.WindowSizeMsg:

		return m, nil
	}

	var cmd tea.Cmd
	m.ta, cmd = m.ta.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {

	statusBar := fmt.Sprintf("%s %v", m.path, m.saved)
	return tea.NewView(statusBar + "\n" + m.ta.View())
}

func (m Model) save() error {
	//0644 -> puedo leer y escribir, los demás solo leer
	return os.WriteFile(m.path, []byte(m.ta.Value()), 0644)
}
