package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			PaddingBottom(1)

	infoStyle = lipgloss.NewStyle().
			PaddingTop(1).
			Faint(true)
)

func (m model) View() string {
	if m.user == nil {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("User not found."),
			infoStyle.Render("Press 'q' to quit"),
		)
	}

	body := fmt.Sprintf(
		"Username: %s\nStatus:   %s",
		m.user.Username,
		m.user.Status,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render("User Profile"),
		lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(0, 1).BorderForeground(lipgloss.Color(m.user.Colour)).Render(body),
		infoStyle.Render("Press 'q' to quit"),
	)
}
