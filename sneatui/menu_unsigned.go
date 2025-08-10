package sneatui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type menuUnassigned struct {
	list list.Model
}

func newMenuUnassigned() tea.Model {
	items := []list.Item{
		item{title: "Sign-in", desc: "Authorize to get access to your data"},
		item{title: "About", desc: "Learn about the Sneat.app"},
	}
	m := menuUnassigned{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Sneat.app - main menu"
	return m
}

func (m menuUnassigned) Init() tea.Cmd {
	return nil
}

func (m menuUnassigned) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m menuUnassigned) View() string {
	return docStyle.Render(m.list.View())
}
