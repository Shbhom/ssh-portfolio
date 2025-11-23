package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "↑"
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "↓"
		case key.Matches(msg, m.keys.Left):
			m.lastKey = "←"
		case key.Matches(msg, m.keys.Right):
			m.lastKey = "→"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
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
