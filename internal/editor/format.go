package editor

import (
	"strings"

	"github.com/juanpablo0606/ttypewrtr/internal/commands"
)

// aplica una transformación a la línea del cursor (títulos, citas, listas)
func (m *Model) applyLine(fn func(string) string) {
	row := m.ta.Line()
	lines := strings.Split(m.ta.Value(), "\n")
	oldLine := lines[row]

	newLine := fn(oldLine)

	diff := len([]rune(newLine)) - len([]rune(oldLine))
	newCol := m.ta.Column() + diff
	if newCol < 0 {
		newCol = 0
	}

	m.replaceCurrentLine(newLine, newCol)
}

// aplica una transformación que también mueve el cursor (negrita, cursiva)
func (m *Model) applyWord(marker string) {
	row := m.ta.Line()
	lines := strings.Split(m.ta.Value(), "\n")
	oldLine := lines[row]

	newLine, newCol := commands.WrapWord(oldLine, m.ta.Column(), marker)

	m.replaceCurrentLine(newLine, newCol)
}

func (m *Model) setCursor(row, col int) {

	if row > m.ta.LineCount()-1 {
		row = m.ta.LineCount() - 1
	}

	m.ta.MoveToBegin()
	for m.ta.Line() < row {
		m.ta.CursorDown()
	}
	m.ta.SetCursorColumn(col)
}

func (m *Model) replaceCurrentLine(newLine string, newCol int) {
	row := m.ta.Line()
	lines := strings.Split(m.ta.Value(), "\n")
	lines[row] = newLine
	m.ta.SetValue(strings.Join(lines, "\n"))

	m.setCursor(row, newCol)
}
