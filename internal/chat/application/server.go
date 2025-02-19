package application

import (
	"context"
	"fmt"
	"github.com/ivofreitas/chat/internal/chat/adapter/websocket"
	"github.com/ivofreitas/chat/pkg/config"
	"github.com/ivofreitas/chat/pkg/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	echo   *echo.Echo
	logger *logrus.Entry
	signal chan struct{}
}

func NewServer() *Server {
	return &Server{
		logger: log.NewEntry(),
		signal: make(chan struct{}),
	}
}

func (s *Server) Run() {
	s.start()
	s.logger.Println("Server started and waiting for the graceful signal...")
	s.watchStop()
}

func (s *Server) start() {
	s.echo = echo.New()
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	go websocket.GetHub().Run()

	register(s.echo)

	env := config.GetEnv()

	s.logger.Infof("Server is starting in port %s.", env.Server.ChatPort)

	addr := fmt.Sprintf(":%s", env.Server.ChatPort)
	go func() {
		if err := s.echo.Start(addr); err != nil {
			s.logger.WithError(err).Fatal("Shutting down the server now")
		}
	}()
}

func (s *Server) watchStop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.logger.Info("Server is stopping...")

	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.logger.Errorln(err)
	}

	close(s.signal)
}
