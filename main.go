package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	wishtea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

type tickMsg time.Time

const (
	padding  = 2
	maxWidth = 80
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
	progress   progress.Model
	prog_perc  float64
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
	loading    bool // will be true until progress bar completes
}

func newModel(userName string) model {
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	return model{
		keys:       keys,
		progress:   prog,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		loading:    true,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
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
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		m.prog_perc += 0.25
		if m.prog_perc > 1.0 {
			m.prog_perc = 1.0
			m.loading = false // 👈 progress finished
			return m, nil
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	// Phase 1: only progress animation
	if m.loading {
		return m.progress.ViewAs(m.prog_perc) + "\n"
	}

	// Phase 2: normal content + help
	helpView := m.help.View(m.keys)

	return m.progress.ViewAs(1.0) + fmt.Sprintf(
		"\nWelcome to your terminal.about.me stub, %s!\n\n"+
			"This is a Bubble Tea TUI running over Wish.\n\n"+
			"Press 'q' to exit.\n",
		m.username,
	) + helpView
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
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
