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
	pauseTicks = 14 // ~1 second at 70ms per tick
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
	username   string
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
	loading    bool // will be true until progress bar completes
	width      int
	height     int

	// intro animation state
	introText  string // "Shubhom Srivastava"
	typedChars int    // how many characters of introText are visible
	cursorOn   bool   // whether to draw the cursor
	phase      int    // 0 = blink only, 1 = typing, 2 = done
	frameCount int    // counts ticks to control timing
}

func newModel(userName string) model {

	return model{
		username:   userName,
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		loading:    true,
		introText:  "Shubhom Srivastava", // what we’ll type out
		typedChars: 0,
		cursorOn:   true,
		phase:      0, // start in blink-only phase
		frameCount: 0,
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
		if !m.loading {
			// If for some reason we still get ticks after loading, ignore them.
			return m, nil
		}

		m.frameCount++

		switch m.phase {
		case 0:
			// Phase 0: just blink the cursor for a bit (~1.5s)
			// Toggle cursor every 3 frames for a nice blink.
			if m.frameCount%3 == 0 {
				m.cursorOn = !m.cursorOn
			}
			// After ~20 frames, switch to typing phase.
			if m.frameCount >= 14 {
				m.phase = 1
				m.frameCount = 0
			}

		case 1:
			// Phase 1: type out the name and keep blinking the cursor.

			// Blink cursor every 3 frames.
			if m.frameCount%3 == 0 {
				m.cursorOn = !m.cursorOn
			}

			// Every 2 frames, reveal one more character.
			if m.frameCount%2 == 0 && m.typedChars < len(m.introText) {
				m.typedChars++
			}

			// When we're done typing, end the loading phase.
			if m.typedChars >= len(m.introText) {
				m.phase = 2
				m.loading = false
				m.cursorOn = false // hide cursor once done
				return m, nil      // stop scheduling ticks
			}

		case 2:
			// Shouldn’t really be hit while loading is true, but just in case:
			if m.frameCount >= pauseTicks {
				m.loading = false  // finally transition to normal mode
				m.cursorOn = false // optional: hide cursor
				return m, nil      // stop scheduling ticks
			}
			return m, nil
		}

		// Keep the animation going.
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

		// Clamp typedChars just in case
		if m.typedChars < 0 {
			m.typedChars = 0
		}
		if m.typedChars > len(m.introText) {
			m.typedChars = len(m.introText)
		}

		visible := m.introText[:m.typedChars]

		cursor := ""
		if m.cursorOn {
			cursor = "▌" // you can change to "_" or "|" if you prefer
		}

		line := visible + cursor
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
	return tea.Tick(70*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
