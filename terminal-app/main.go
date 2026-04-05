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
	defaultHost        = "0.0.0.0"
	defaultPort        = "22"
	defaultHostKeyPath = ".ssh/id_ed25519"
)

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

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
	pbBaseURL := os.Getenv("PB_BASE_URL")
	if pbBaseURL == "" {
		log.Fatal("PB_BASE_URL environment variable not set")
	}

	host := firstNonEmpty(os.Getenv("SSH_HOST"), defaultHost)
	port := firstNonEmpty(os.Getenv("SSH_PORT"), os.Getenv("PORT"), defaultPort)
	hostKeyPath := firstNonEmpty(os.Getenv("SSH_HOST_KEY_PATH"), defaultHostKeyPath)

	lipgloss.SetColorProfile(termenv.TrueColor)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
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
	log.Info("Starting SSH server", "host", host, "port", port, "pb_base_url", pbBaseURL, "host_key_path", hostKeyPath)
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
