package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// AdaptiveColor pairs for each platform: Light bg → Dark bg.
// Dark-on-dark colours (GitHub #333333, Dev.to #0A0A0A) get light variants
// for dark terminal backgrounds.
var platformColours = map[string]lipgloss.AdaptiveColor{
	"github":    {Light: "#333333", Dark: "#cccccc"},
	"linkedin":  {Light: "#0A66C2", Dark: "#0A66C2"},
	"youtube":   {Light: "#FF0000", Dark: "#FF4444"},
	"twitter":   {Light: "#1DA1F2", Dark: "#1DA1F2"},
	"instagram": {Light: "#E4405F", Dark: "#E4405F"},
	"mastodon":  {Light: "#6364FF", Dark: "#8485FF"},
	"bluesky":   {Light: "#0085FF", Dark: "#38a8ff"},
	"devto":     {Light: "#0A0A0A", Dark: "#ffffff"},
	"hashnode":  {Light: "#2962FF", Dark: "#5c87ff"},
	"discord":   {Light: "#5865F2", Dark: "#7289DA"},
	"twitch":    {Light: "#9146FF", Dark: "#a970ff"},
}

var platformNames = map[string]string{
	"github":    "GitHub",
	"linkedin":  "LinkedIn",
	"youtube":   "YouTube",
	"twitter":   "Twitter",
	"instagram": "Instagram",
	"mastodon":  "Mastodon",
	"bluesky":   "Bluesky",
	"devto":     "Dev.to",
	"hashnode":  "Hashnode",
	"discord":   "Discord",
	"twitch":    "Twitch",
	"website":   "Website",
	"blog":      "Blog",
	"email":     "Email",
}

func (m model) View() string {
	r := m.renderer
	if r == nil {
		r = lipgloss.DefaultRenderer()
	}

	if m.user == nil {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			r.NewStyle().Foreground(lipgloss.Color("9")).Render("User not found."),
			r.NewStyle().PaddingTop(1).Faint(true).Render("Press 'q' to quit"),
		)
	}

	layout := m.user.Layout
	if layout == "" {
		layout = "full"
	}

	// Available URL width after border (2), padding (4), "▸ " (2), label col, padding (2)
	contentWidth := m.width - 6
	if contentWidth < 20 {
		contentWidth = 0
	}

	var lines []string

	// --- Name ---
	name := m.user.Username
	if m.user.DisplayName != "" {
		name = m.user.DisplayName
	}
	lines = append(lines, r.NewStyle().Bold(true).Render(name))

	// --- Meta (full only) ---
	if layout == "full" {
		var metaParts []string
		if m.user.DisplayName != "" {
			metaParts = append(metaParts, "@"+m.user.Username)
		}
		if m.user.Pronouns != "" {
			metaParts = append(metaParts, m.user.Pronouns)
		}
		if m.user.Location != "" {
			metaParts = append(metaParts, m.user.Location)
		}
		if len(metaParts) > 0 {
			lines = append(lines, r.NewStyle().Faint(true).Render(strings.Join(metaParts, " · ")))
		}
	}

	// --- Bio / status ---
	switch layout {
	case "minimal":
		// omit
	case "compact":
		if m.user.Status != "" {
			lines = append(lines, "", m.user.Status)
		}
	default: // full
		if m.user.Bio != "" {
			lines = append(lines, "", m.user.Bio)
		}
		if m.user.Status != "" {
			lines = append(lines, m.user.Status)
		}
	}

	// --- Links ---
	if len(m.links) > 0 {
		lines = append(lines, r.NewStyle().Bold(true).MarginTop(1).Render("Links"))
		lines = append(lines, r.NewStyle().Faint(true).Render("─────"))

		maxLabelLen := 0
		for _, link := range m.links {
			if l := resolveLabel(link); len(l) > maxLabelLen {
				maxLabelLen = len(l)
			}
		}

		urlMaxLen := 0
		if contentWidth > 0 {
			urlMaxLen = contentWidth - maxLabelLen - 4 // "▸ " + label + padding
		}

		for _, link := range m.links {
			label := resolveLabel(link)
			labelStyle := r.NewStyle().Bold(true).Foreground(resolvePlatformColour(link, m.user.Colour))
			urlStyle := r.NewStyle().Faint(true)

			padding := strings.Repeat(" ", maxLabelLen-len(label)+2)
			url := trimScheme(link.URL)
			if urlMaxLen > 3 && len(url) > urlMaxLen {
				url = url[:urlMaxLen-3] + "..."
			}
			lines = append(lines, fmt.Sprintf("▸ %s%s%s", labelStyle.Render(label), padding, urlStyle.Render(url)))
		}
	}

	// --- SSH command ---
	if m.user.Username != "" {
		lines = append(lines, r.NewStyle().Faint(true).MarginTop(1).Render(
			fmt.Sprintf("ssh %s@ssh-me.com", m.user.Username),
		))
	}

	content := strings.Join(lines, "\n")

	bordered := r.NewStyle().
		Border(lipgloss.ThickBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color(m.user.Colour)).
		Render(content)

	brandingStyle := r.NewStyle().Faint(true).Italic(true)
	quitHint := r.NewStyle().Faint(true)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		bordered,
		brandingStyle.Render("powered by ssh-me"),
		quitHint.Render("Press 'q' to quit"),
	)
}

func resolveLabel(link Link) string {
	if link.Label != "" {
		return link.Label
	}
	if name, ok := platformNames[link.Platform]; ok {
		return name
	}
	return "Link"
}

// resolvePlatformColour returns an AdaptiveColor for known platforms,
// falling back to the user's own colour for website/blog/email/custom.
func resolvePlatformColour(link Link, userColour string) lipgloss.TerminalColor {
	if colour, ok := platformColours[link.Platform]; ok {
		return colour
	}
	if userColour != "" {
		return lipgloss.Color(userColour)
	}
	return lipgloss.Color("12") // bright blue fallback
}

func trimScheme(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	return url
}
