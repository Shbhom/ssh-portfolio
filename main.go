package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	wishtea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
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

type model struct {
	username string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf(
		"Welcome to your terminal.about.me stub, %s!\n\n"+
			"This is a Bubble Tea TUI running over Wish.\n\n"+
			"Press 'q' to exit.\n",
		m.username,
	)
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := model{
		username: s.User(),
	}

	opts := []tea.ProgramOption{
		tea.WithInput(s),
		tea.WithOutput(s),
		tea.WithAltScreen(), // optional but nice
	}

	return m, opts
}
