package ui

import "github.com/charmbracelet/lipgloss"

var (
	nameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD7FF"))

	// Thick, colored cursor
	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87"))
)
