package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

// Ripped from docs:
// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	var user *Profile
	if v := s.Context().Value(userCtxKey); v != nil {
		if u, ok := v.(*Profile); ok {
			user = u
		}
	}

	var links []Link
	if user != nil {
		if l, err := GetLinksByProfileId(os.Getenv("PB_BASE_URL"), user.Id); err == nil {
			links = l
		}
	}

	// Ripped from docs:
	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	m := model{
		term:      pty.Term,
		profile:   renderer.ColorProfile().Name(),
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		renderer:  renderer,
		txtStyle:  txtStyle,
		quitStyle: quitStyle,
		user:      user,
		links:     links,
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
