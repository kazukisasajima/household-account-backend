package web

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"household-account-backend/adapter/controller/echo/router"
)

type EchoServer struct {
	router     *echo.Echo
	host, port string
}

func NewEchoServer(host, port string, db *gorm.DB) (Server, error) {
	return &EchoServer{
		router: router.NewEchoRouter(db),
		host:   host,
		port:   port,
	}, nil
}

func (e *EchoServer) Start() error {
	return e.router.Start(fmt.Sprintf("%s:%s", e.host, e.port))
}

func (e *EchoServer) Shutdown(ctx context.Context) error {
	return e.router.Shutdown(ctx)
}
