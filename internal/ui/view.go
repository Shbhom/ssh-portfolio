package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var content string

	if m.loading {
		// Intro phase: typewriter animation for "Shubhom Srivastava"

		if m.typedChars < 0 {
			m.typedChars = 0
		}
		if m.typedChars > len(m.introText) {
			m.typedChars = len(m.introText)
		}

		visible := m.introText[:m.typedChars]

		// Thicker cursor: full block "█" (or "▋"/"▌" if you want slimmer)
		cursorChar := ""
		if m.cursorOn {
			cursorChar = "█"
		}

		// Style the name + cursor separately
		styledName := m.nameStyle.Render(visible)
		styledCursor := ""
		if cursorChar != "" {
			styledCursor = m.cursorStyle.Render(cursorChar)
		}

		line := styledName + styledCursor
		content = line
	} else {
		// Normal phase: main content + help
		helpView := m.help.View(m.keys)

		body := fmt.Sprintf(
			"Welcome to your terminal.about.me stub, %s!\n\n"+
				"This is a Bubble Tea TUI running over Wish.\n\n"+
				"Press 'q' to exit.\n",
			m.username,
		)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			body,
			"",
			helpView,
		)
	}

	// If we don't know the size yet, just return the raw content.
	if m.width == 0 || m.height == 0 {
		return content + "\n"
	}

	// Center the content in the available space
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center, // horizontal
		lipgloss.Center, // vertical
		content,
	)
}
