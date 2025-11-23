package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	pauseTicks = 10 // ~1 second at 70ms per tick
)

type tickMsg time.Time

type model struct {
	username string
	keys     keyMap
	help     help.Model
	lastKey  string
	quitting bool
	loading  bool // will be true until progress bar completes
	width    int
	height   int

	// intro animation state
	introText  string // "Shubhom Srivastava"
	typedChars int    // how many characters of introText are visible
	cursorOn   bool   // whether to draw the cursor
	phase      int    // 0 = blink only, 1 = typing, 2 = done
	frameCount int    // counts ticks to control timing

	// styles
	nameStyle   lipgloss.Style
	cursorStyle lipgloss.Style
}

func NewModel(userName string) model {

	return model{
		username:   userName,
		keys:       keys,
		help:       help.New(),
		loading:    true,
		introText:  "Shubhom Srivastava", // what we’ll type out
		typedChars: 0,
		cursorOn:   true,
		phase:      0, // start in blink-only phase
		frameCount: 0,
		// Bold, slightly “bigger-feeling” name
		nameStyle: nameStyle, // tweak color if you want

		// Thick, colored cursor
		cursorStyle: cursorStyle,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(), // your progress timer
		tea.SetWindowTitle("Shubhom's Portfolio"), // 👈 window title
	)
}
