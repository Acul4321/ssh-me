package main

import "github.com/charmbracelet/lipgloss"

type model struct {
	term      string
	profile   string
	width     int
	height    int
	bg        string
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
	user      *User
}
