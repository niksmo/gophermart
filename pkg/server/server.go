package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

const shutdownTimeout = 5 * time.Second

type HTTPServer struct {
	*fiber.App
	addr   string
	logger zerolog.Logger
}

func NewHTTPServer(addr string, logger zerolog.Logger) HTTPServer {
	return HTTPServer{App: fiber.New(), addr: addr, logger: logger}
}

func (s HTTPServer) Run() {
	if err := s.Listen(s.addr); err != nil {
		s.logger.Info().Err(err).Msg("server listening error")
	}
}

func (s HTTPServer) Close() {
	if err := s.ShutdownWithTimeout(shutdownTimeout); err != nil {
		s.logger.Warn().Err(err).Msg("closing server connections")
	} else {
		s.logger.Info().Msg("server connection safely closed")
	}
}
