package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.activeTab == 1 { // assuming 0=Overview, 1=Experience
			switch msg.String() {
			case "j", "down":
				m.expList.NextPage()
			case "k", "up":
				m.expList.PrevPage()
			}
		}
		if m.activeTab == 2 { // assuming 0=Overview, 1=Experience
			switch msg.String() {
			case "j", "down":
				m.projList.NextPage()
			case "k", "up":
				m.projList.PrevPage()
			}
		}
		switch msg.String() {
		case "1":
			m.activeTab = 0
		case "2":
			m.activeTab = 1
		case "3":
			m.activeTab = 2
		case "4":
			m.activeTab = 3

		case "left", "h":
			if m.activeTab == 0 {
				m.activeTab = 3
			} else {
				m.activeTab--
			}
		case "right", "l":
			if m.activeTab == 3 {
				m.activeTab = 0
			} else {
				m.activeTab++
			}
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width // 👈 store
		m.height = msg.Height
		return m, nil

	case tickMsg:
		// If loading is already over, ignore further ticks
		if !m.loading {
			return m, nil
		}

		m.frameCount++

		switch m.phase {
		case 0:
			// Phase 0: blink-only cursor for a bit (~1.5s)
			if m.frameCount%3 == 0 {
				m.cursorOn = !m.cursorOn
			}
			if m.frameCount >= 20 {
				m.phase = 1
				m.frameCount = 0
			}

		case 1:
			// Phase 1: typing the name
			if m.frameCount%3 == 0 {
				m.cursorOn = !m.cursorOn
			}

			if m.frameCount%2 == 0 && m.typedChars < len(m.introText) {
				m.typedChars++
			}

			if m.typedChars >= len(m.introText) {
				// ✅ DO *NOT* set loading = false here
				// just move to pause phase
				m.phase = 2
				m.frameCount = 0
				m.cursorOn = true // keep cursor visible at end of name
			}

		case 2:
			// Phase 2: pause with full name visible
			if m.frameCount >= pauseTicks {
				// Now we finally end loading and switch to normal mode
				m.loading = false
				m.cursorOn = false // optional: hide cursor in normal mode
				return m, nil      // stop scheduling ticks
			}
		}

		// Keep animation running for phases 0, 1, and 2
		return m, tickCmd()
	}
	return m, nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
