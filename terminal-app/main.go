package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

const userCtxKey = "pbProfile"

const (
	host = "0.0.0.0"
	port = "22"
)

func userLookupMiddleware(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		username := sess.User()
		profile, err := GetProfileByUsername(os.Getenv("PB_BASE_URL"), username)
		if err != nil {
			log.Warn("PB lookup failed", "user", username, "err", err)
		} else {
			sess.Context().SetValue(userCtxKey, profile)
		}
		next(sess)
	}
}

func main() {
	if os.Getenv("PB_BASE_URL") == "" {
		log.Fatal("PB_BASE_URL environment variable not set")
	}

	lipgloss.SetColorProfile(termenv.TrueColor)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			userLookupMiddleware,
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
