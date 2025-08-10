package sneatui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Navigation messages for app-level routing.
type (
	navToLoginMsg        struct{}
	navToAboutMsg        struct{}
	navBackToUnsignedMsg struct{}
	navSignedInMsg       struct{}
	navSignOutMsg        struct{}
)

// appModel is the root model that owns application state and active screen.
type appModel struct {
	active string // "unsigned", "login", "about", "signed"

	unsigned menuUnsigned
	signed   menuSignedIn
	login    loginModel
	about    aboutModel

	winW int
	winH int
}

func newAppModel() tea.Model {
	// Initialize child models with defaults
	mu := newMenuUnassigned().(menuUnsigned)
	return appModel{
		active:   "unsigned",
		unsigned: mu,
	}
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winW, m.winH = msg.Width, msg.Height
		// Forward size to active child
		return m.forwardToActive(msg)
	case navToLoginMsg:
		// create login model fresh
		m.login = newLoginModel().(loginModel)
		m.active = "login"
		// Propagate last window size to child if we have it
		if m.winW > 0 && m.winH > 0 {
			return m, func() tea.Msg { return tea.WindowSizeMsg{Width: m.winW, Height: m.winH} }
		}
		return m, nil
	case navToAboutMsg:
		m.about = newAboutModel().(aboutModel)
		m.active = "about"
		if m.winW > 0 && m.winH > 0 {
			return m, func() tea.Msg { return tea.WindowSizeMsg{Width: m.winW, Height: m.winH} }
		}
		return m, nil
	case navBackToUnsignedMsg:
		m.active = "unsigned"
		if m.winW > 0 && m.winH > 0 {
			return m, func() tea.Msg { return tea.WindowSizeMsg{Width: m.winW, Height: m.winH} }
		}
		return m, nil
	case navSignedInMsg:
		m.signed = newMenuSignedIn().(menuSignedIn)
		m.active = "signed"
		if m.winW > 0 && m.winH > 0 {
			return m, func() tea.Msg { return tea.WindowSizeMsg{Width: m.winW, Height: m.winH} }
		}
		return m, nil
	case navSignOutMsg:
		m.active = "unsigned"
		if m.winW > 0 && m.winH > 0 {
			return m, func() tea.Msg { return tea.WindowSizeMsg{Width: m.winW, Height: m.winH} }
		}
		return m, nil
	default:
		return m.forwardToActive(msg)
	}
}

func (m appModel) forwardToActive(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.active {
	case "login":
		child, cmd := m.login.Update(msg)
		m.login = child.(loginModel)
		return m, cmd
	case "about":
		child, cmd := m.about.Update(msg)
		m.about = child.(aboutModel)
		return m, cmd
	case "signed":
		child, cmd := m.signed.Update(msg)
		m.signed = child.(menuSignedIn)
		return m, cmd
	default: // "unsigned"
		child, cmd := m.unsigned.Update(msg)
		m.unsigned = child.(menuUnsigned)
		return m, cmd
	}
}

func (m appModel) View() string {
	switch m.active {
	case "login":
		return m.login.View()
	case "about":
		return m.about.View()
	case "signed":
		return m.signed.View()
	default:
		return m.unsigned.View()
	}
}
