package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/textarea"
	"github.com/juanpablo0606/ttypewrtr/internal/editor"
)

type app struct {
	ta textarea.Model
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing arguments.\nUse: ttypewrtr <file.md>")
		os.Exit(1)
	}
	path := os.Args[1]

	//initialize bubbletea
	if _, err := editor.New(editor.Model{}, path); err != nil {
		log.Fatal(err)
	}

}
