package ui

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
)

var (
	nameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD7FF"))

	// Thick, colored cursor
	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87"))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2).
			Width(appWidth).
			Height(appHeight)

	contentStyle = lipgloss.NewStyle().
			Width(appWidth - 4). // inside padding
			Height(appHeight - 5)

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFD7FF")).
			Padding(0, 1)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Padding(0, 1)

	tabsRowStyle = lipgloss.NewStyle().
			Width(appWidth).
			Align(lipgloss.Center)

	footerStyle = lipgloss.NewStyle().
			Width(appWidth).
			Foreground(lipgloss.Color("#555555")).
			Align(lipgloss.Center)
	contactTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFD7FF")). // soft pink, same vibe as tabs
				MarginBottom(1).
				Underline(false)
)

func newPaginator(total int) paginator.Model {
	p := paginator.New()
	p.PerPage = 1
	p.TotalPages = total
	p.Type = paginator.Dots

	active := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF75B7")). // pink
		Bold(true).
		Render("●")

	inactive := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")). // dim gray
		Render("•")

	p.ActiveDot = active
	p.InactiveDot = inactive

	return p
}
