package sneatui

import tea "github.com/charmbracelet/bubbletea"

func InitialModel() tea.Model {
	return newMenuUnassigned()
}
