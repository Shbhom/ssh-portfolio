package ui

import (
	"time"

	"github.com/Shbhom/ssh-portfolio/internal/portfolio"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	pauseTicks = 10 // ~1 second at 70ms per tick
	appWidth   = 100
	appHeight  = 20
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

	activeTab int // 0 = Overview, 1 = Experience, 2 = Projects, 3 = Contact
	portfolio *portfolio.Portfolio
	expList   paginator.Model
}

func NewModel(userName string, p *portfolio.Portfolio) model {

	// expPager :=
	// expPager.Type = paginator.Dots
	// expPager.PerPage = 1
	// if len(p.Experiences) == 0 {
	// 	expPager.SetTotalPages(1)
	// } else {
	// 	expPager.SetTotalPages(len(p.Experiences))
	// }

	expPager := newPaginator(len(p.Experiences))

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
		activeTab:  0,

		nameStyle: nameStyle,

		cursorStyle: cursorStyle,
		portfolio:   p,
		expList:     expPager,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(), // your progress timer
		tea.SetWindowTitle("Shubhom's Portfolio"), // 👈 window title
	)
}
