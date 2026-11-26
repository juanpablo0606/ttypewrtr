package doc

import "strings"

type Style uint8

const (
	Heading Style = 1 << iota
	Bold
	Italic
	Quote
	Code
	Marker
)

var styleNames = []string{"Heading", "Bold", "Italic", "Quote", "Code", "Marker"}

const Plain Style = 0

func (s Style) Has(other Style) bool {
	return s&other != 0
}

func (s Style) String() string {
	nombres := []string{}
	if s == Plain {
		return "Plain"
	}
	for i, name := range styleNames {
		if s.Has(Style(1 << i)) {
			nombres = append(nombres, name)
		}
	}
	return strings.Join(nombres, "|")
}
