package sneatui

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// aboutModel shows the content of a markdown file describing Sneat.app.
// ESC should return to the unsigned menu. We keep rendering simple (no markdown parsing)
// to satisfy minimal-diff requirement.
type aboutModel struct {
	content string
	style   lipgloss.Style
	winW    int
	winH    int
	loaded  bool
}

func newAboutModel() tea.Model {
	return aboutModel{style: lipgloss.NewStyle().Margin(1, 2)}
}

func (m aboutModel) Init() tea.Cmd {
	return func() tea.Msg {
		// Load README.md from project root. Tests run with CWD at package dir; go test
		// executes package in its directory, so root README is one level up.
		// If load fails, we still show a friendly message.
		path := filepath.Clean("../README.md")
		b, err := os.ReadFile(path)
		if err != nil {
			return aboutLoadedMsg("Sneat.app\n\nAbout information is unavailable: " + err.Error())
		}
		return aboutLoadedMsg(string(b))
	}
}

type aboutLoadedMsg string

func (m aboutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Ask app to go back to unsigned menu
			return m, func() tea.Msg { return navBackToUnsignedMsg{} }
		}
	case tea.WindowSizeMsg:
		m.winW, m.winH = msg.Width, msg.Height
		// nothing else needed; we use style margins in View
	case aboutLoadedMsg:
		m.content = string(msg)
		m.loaded = true
	}
	return m, nil
}

func (m aboutModel) View() string {
	// Provide a header; keep simple
	header := lipgloss.NewStyle().Bold(true).Render("About Sneat.app")
	body := m.content
	if body == "" {
		body = "Loading…"
	}
	return m.style.Render(header + "\n\n" + body + "\n\n(ESC to return)")
}
