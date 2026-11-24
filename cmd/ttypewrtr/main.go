package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/juanpablo0606/ttypewrtr/internal/editor"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing arguments.\nUse: ttypewrtr <file.md>")
		os.Exit(1)
	}
	path := os.Args[1]

	m, err := editor.New(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//initialize bubbletea
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
