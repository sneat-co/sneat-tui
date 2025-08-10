package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sneat-co/sneat-tui/sneatui"
	"os"
)

// test hooks to allow overriding in tests
var (
	getProgram = newProgram
	exit       = os.Exit
)

func main() {
	p := getProgram()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		exit(1)
	}
}

type program interface {
	Run() (tea.Model, error)
}

func newProgram() program {
	return tea.NewProgram(sneatui.InitialModel(), tea.WithAltScreen())
}
