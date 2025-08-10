package sneatui

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
)

func TestInitialModel_ReturnsAppModel_ShowsUnsignedMenu(t *testing.T) {
	m := InitialModel()
	am, ok := m.(appModel)
	if !ok {
		t.Fatalf("InitialModel returned %T, want appModel", m)
	}
	view := am.View()
	if !strings.Contains(view, "Sneat.app - main menu") {
		t.Fatalf("initial view missing main menu title; view=\n%s", view)
	}
}

func TestApp_Unsigned_CtrlCQuits(t *testing.T) {
	m := InitialModel().(appModel)
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := model.(appModel); !ok {
		t.Fatalf("Update returned model %T, want appModel", model)
	}
	if cmd == nil {
		t.Fatalf("cmd is nil, want tea.Quit")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("cmd() returned %T, want tea.QuitMsg", msg)
	}
}

func TestApp_Unsigned_WindowSizeSetsListSize(t *testing.T) {
	m := InitialModel().(appModel)
	// initial sizes
	w0, h0 := m.unsigned.list.Width(), m.unsigned.list.Height()

	hMargin, vMargin := docStyle.GetFrameSize()
	msg := tea.WindowSizeMsg{Width: 100, Height: 40}
	model, _ := m.Update(msg)
	am := model.(appModel)
	wantW := msg.Width - hMargin
	wantH := msg.Height - vMargin
	if am.unsigned.list.Width() != wantW || am.unsigned.list.Height() != wantH {
		t.Fatalf("list size = (%d,%d), want (%d,%d); initial was (%d,%d)", am.unsigned.list.Width(), am.unsigned.list.Height(), wantW, wantH, w0, h0)
	}
}
