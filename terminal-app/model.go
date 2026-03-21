package main

import "github.com/charmbracelet/lipgloss"

type model struct {
	term      string
	profile   string
	width     int
	height    int
	renderer  *lipgloss.Renderer
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
	user      *Profile
	links     []Link
}
