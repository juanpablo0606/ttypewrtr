package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/juanpablo0606/ttypewrtr/internal/config"
	"github.com/juanpablo0606/ttypewrtr/internal/editor"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing arguments.\nUse: ttypewrtr <file.md>")
		os.Exit(1)
	}
	path := os.Args[1]

	cfgPath, err := config.Path()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, err := editor.New(path, cfg.Keys)
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
