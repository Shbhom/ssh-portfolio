package sshserver

import (
	"github.com/Shbhom/ssh-portfolio/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	wishtea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func New(addr, hostKeyPath string) (*ssh.Server, error) {
	srv, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(hostKeyPath),
		ssh.AllocatePty(),
		wish.WithMiddleware(
			logging.Middleware(),
			wishtea.Middleware(teaHandler),
		),
	)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := ui.NewModel(s.User())

	opts := []tea.ProgramOption{
		tea.WithInput(s),
		tea.WithOutput(s),
		tea.WithAltScreen(), // optional but nice
	}

	return m, opts
}
