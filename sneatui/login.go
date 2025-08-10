package sneatui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// loginModel presents a simple login form with email and password inputs.
// ESC returns to the main menu. This model intentionally does not perform
// any real authentication; it exists to satisfy navigation requirements.
type loginModel struct {
	styles   loginStyles
	email    textinput.Model
	password textinput.Model
	focused  int // 0=email, 1=password
	winW     int
	winH     int
}

type loginStyles struct {
	doc    lipgloss.Style
	header lipgloss.Style
	label  lipgloss.Style
}

func newLoginModel() tea.Model {
	// Email input
	e := textinput.New()
	e.Placeholder = "you@example.com"
	e.Focus()
	e.Prompt = ""
	e.CharLimit = 256
	e.Width = 40

	// Password input
	p := textinput.New()
	p.Placeholder = "password"
	p.Prompt = ""
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.CharLimit = 256
	p.Width = 40

	st := loginStyles{
		doc:    lipgloss.NewStyle().Margin(1, 2),
		header: lipgloss.NewStyle().Bold(true),
		label:  lipgloss.NewStyle().Faint(true),
	}

	return loginModel{styles: st, email: e, password: p, focused: 0}
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Ask app to go back to unsigned menu
			return m, func() tea.Msg { return navBackToUnsignedMsg{} }
		case "tab", "shift+tab":
			if m.focused == 0 {
				m.focused = 1
				m.password.Focus()
				m.email.Blur()
			} else {
				m.focused = 0
				m.email.Focus()
				m.password.Blur()
			}
			return m, nil
		case "enter":
			// Move focus from email to password on enter; if already on password, signal sign-in
			if m.focused == 0 {
				m.focused = 1
				m.password.Focus()
				m.email.Blur()
				return m, nil
			}
			// When on password, treat Enter as successful sign-in for navigation purposes
			return m, func() tea.Msg { return navSignedInMsg{} }
		}
	case tea.WindowSizeMsg:
		// Save the window size and adjust input widths to available space minus margins
		m.winW, m.winH = msg.Width, msg.Height
		h, _ := m.styles.doc.GetFrameSize()
		w := msg.Width - h
		if w < 20 {
			w = 20
		}
		m.email.Width = w
		m.password.Width = w
	}

	var cmd tea.Cmd
	if m.focused == 0 {
		m.email, cmd = m.email.Update(msg)
		return m, cmd
	}
	m.password, cmd = m.password.Update(msg)
	return m, cmd
}

func (m loginModel) View() string {
	header := m.styles.header.Render("Sneat.app - Sign in")
	emailLbl := m.styles.label.Render("Email:")
	passLbl := m.styles.label.Render("Password:")
	instructions := m.styles.label.Render("TAB to switch, ESC to return")
	view := header + "\n\n" + emailLbl + "\n" + m.email.View() + "\n\n" + passLbl + "\n" + m.password.View() + "\n\n" + instructions
	return m.styles.doc.Render(view)
}
