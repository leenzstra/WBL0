package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

type Server struct {
	notify chan error
	logger *zap.Logger
	*fiber.App
}

func New(app *fiber.App, logger *zap.Logger) *Server {
	s := &Server{
		notify: make(chan error, 1),
		App:    app,
		logger: logger,
	}

	logger.Info("Http server started")

	return s
}

func (s *Server) Listen(addr string) {
	go func() {
		s.notify <- s.App.Listen(addr)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.App.ShutdownWithTimeout(shutdownTimeout)
}
