package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) viewTabs() string {
	labels := []string{"Overview", "Experience", "Projects", "Contact"}
	var rendered []string

	for i, label := range labels {
		if i == m.activeTab {
			rendered = append(rendered, tabActiveStyle.Render(label))
		} else {
			rendered = append(rendered, tabInactiveStyle.Render(label))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, rendered...)
	return tabsRowStyle.Render(row)
}

func (m model) viewTabContent() string {
	var text string

	switch m.activeTab {
	case 0:
		text = fmt.Sprintf("%s\n%s\n%s\n", m.portfolio.Name, m.portfolio.Tagline, m.portfolio.Overview.Intro)
	case 1:
		text = "Experience tab — placeholder content."
	case 2:
		text = "Projects tab — placeholder content."
	case 3:
		text = "Contact tab — placeholder content."
	}

	return contentStyle.Render(text)
}

func (m model) viewFooter() string {
	helpLine := "h/← & l/→: switch tabs  •  1–4: jump to tab  •  q: quit"
	return footerStyle.Render(helpLine)
}

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

		// Thicker cursor: full block "█"
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
		// 🔹 Main portfolio card view

		tabContent := m.viewTabContent()
		tabsRow := m.viewTabs()
		footer := m.viewFooter()

		body := lipgloss.JoinVertical(
			lipgloss.Left,
			tabContent,
			"",
			tabsRow,
			footer,
		)

		content = cardStyle.Render(body)
	}

	// If we don't know the size yet, just return the raw content.
	if m.width == 0 || m.height == 0 {
		return content + "\n"
	}

	// Center the content (intro or card) in the available space
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center, // horizontal
		lipgloss.Center, // vertical
		content,
	)
}
