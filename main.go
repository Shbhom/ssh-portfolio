package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	wishtea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

type tickMsg time.Time

const (
	pauseTicks = 10 // ~1 second at 70ms per tick
)

func main() {
	// Our SSH app will listen on all interfaces, port 23234
	addr := ":23234"

	// Host key will be created automatically at this path if it doesn't exist.
	hostKeyPath := "ssh_host_ed25519"

	srv, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(hostKeyPath),

		// NEW: ensure a PTY is allocated for interactive sessions
		ssh.AllocatePty(),

		wish.WithMiddleware(
			logging.Middleware(),
			wishtea.Middleware(teaHandler),
		),
	)

	if err != nil {
		log.Fatalf("failed to create wish server: %v", err)
	}

	log.Printf("Starting SSH server on %s ...", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

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

func newModel(userName string) model {

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
		nameStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD7FF")), // tweak color if you want

		// Thick, colored cursor
		cursorStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87")),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(), // your progress timer
		tea.SetWindowTitle("Shubhom's Portfolio"), // 👈 window title
	)
}

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

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := newModel(s.User())

	opts := []tea.ProgramOption{
		tea.WithInput(s),
		tea.WithOutput(s),
		tea.WithAltScreen(), // optional but nice
	}

	return m, opts
}

func tickCmd() tea.Cmd {
	return tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
